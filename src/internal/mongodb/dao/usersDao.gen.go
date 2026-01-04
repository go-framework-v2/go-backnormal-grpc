package dao

import (
	"context"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/model/bo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UsersDao MongoDB 数据访问接口
type UsersDao interface {
	// 基础 CRUD 操作
	FindByID(id primitive.ObjectID) (*bo.UsersBo, error)
	FindOne(filter bson.M) (*bo.UsersBo, error)
	Find(filter bson.M, opts ...*options.FindOptions) ([]*bo.UsersBo, error)
	Insert(obj *bo.UsersBo) (*bo.UsersBo, error)
	UpdateByID(id primitive.ObjectID, update bson.M) (*bo.UsersBo, error)
	DeleteByID(id primitive.ObjectID) error
	Count(filter bson.M) (int64, error)

	// 事务支持
	BeginTransaction() (TransactionSession, error)
	WithTransaction(fn func(txSession TransactionSession) error) error
}

// TransactionSession 事务会话接口
type TransactionSession interface {
	CommitTransaction() error
	AbortTransaction() error
	EndSession()

	FindByID(id primitive.ObjectID) (*bo.UsersBo, error)
	FindOne(filter bson.M) (*bo.UsersBo, error)
	Insert(obj *bo.UsersBo) (*bo.UsersBo, error)
	UpdateByID(id primitive.ObjectID, update bson.M) (*bo.UsersBo, error)
	DeleteByID(id primitive.ObjectID) error
	Count(filter bson.M) (int64, error)
}

// transactionSessionImpl 事务会话实现
type transactionSessionImpl struct {
	session    mongo.Session
	collection *mongo.Collection
}

// customUsersDao 自定义 DAO 实现
type customUsersDao struct {
	collection *mongo.Collection
	db         *mongo.Database
	client     *mongo.Client
}

// NewUsersDao 创建新的 DAO 实例
func NewUsersDao(db *mongo.Database) UsersDao {
	collection := db.Collection("users")
	client := db.Client()
	return &customUsersDao{
		collection: collection,
		db:         db,
		client:     client,
	}
}

// 基础 CRUD 方法
func (d *customUsersDao) FindByID(id primitive.ObjectID) (*bo.UsersBo, error) {
	ctx := context.Background()
	var result bo.UsersBo
	err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (d *customUsersDao) FindOne(filter bson.M) (*bo.UsersBo, error) {
	ctx := context.Background()
	var result bo.UsersBo
	err := d.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (d *customUsersDao) Find(filter bson.M, opts ...*options.FindOptions) ([]*bo.UsersBo, error) {
	ctx := context.Background()
	cursor, err := d.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*bo.UsersBo
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *customUsersDao) Insert(obj *bo.UsersBo) (*bo.UsersBo, error) {
	ctx := context.Background()

	if obj.ID.IsZero() {
		obj.ID = primitive.NewObjectID()
	}

	_, err := d.collection.InsertOne(ctx, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (d *customUsersDao) UpdateByID(id primitive.ObjectID, update bson.M) (*bo.UsersBo, error) {
	ctx := context.Background()
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var result bo.UsersBo
	err := d.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *customUsersDao) DeleteByID(id primitive.ObjectID) error {
	ctx := context.Background()
	_, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (d *customUsersDao) Count(filter bson.M) (int64, error) {
	ctx := context.Background()
	return d.collection.CountDocuments(ctx, filter)
}

// 事务支持
func (d *customUsersDao) BeginTransaction() (TransactionSession, error) {
	ctx := context.Background()
	session, err := d.client.StartSession()
	if err != nil {
		return nil, err
	}

	err = session.StartTransaction()
	if err != nil {
		session.EndSession(ctx)
		return nil, err
	}

	return &transactionSessionImpl{
		session:    session,
		collection: d.collection,
	}, nil
}

func (d *customUsersDao) WithTransaction(fn func(txSession TransactionSession) error) error {
	txSession, err := d.BeginTransaction()
	if err != nil {
		return err
	}
	defer txSession.EndSession()

	err = fn(txSession)
	if err != nil {
		if abortErr := txSession.AbortTransaction(); abortErr != nil {
			return abortErr
		}
		return err
	}

	return txSession.CommitTransaction()
}

// 事务会话实现
func (ts *transactionSessionImpl) CommitTransaction() error {
	ctx := context.Background()
	return ts.session.CommitTransaction(ctx)
}

func (ts *transactionSessionImpl) AbortTransaction() error {
	ctx := context.Background()
	return ts.session.AbortTransaction(ctx)
}

func (ts *transactionSessionImpl) EndSession() {
	ctx := context.Background()
	ts.session.EndSession(ctx)
}

func (ts *transactionSessionImpl) FindByID(id primitive.ObjectID) (*bo.UsersBo, error) {
	ctx := context.Background()
	var result bo.UsersBo
	err := mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		return ts.collection.FindOne(sc, bson.M{"_id": id}).Decode(&result)
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (ts *transactionSessionImpl) FindOne(filter bson.M) (*bo.UsersBo, error) {
	ctx := context.Background()
	var result bo.UsersBo
	err := mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		return ts.collection.FindOne(sc, filter).Decode(&result)
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (ts *transactionSessionImpl) Insert(obj *bo.UsersBo) (*bo.UsersBo, error) {
	ctx := context.Background()

	if obj.ID.IsZero() {
		obj.ID = primitive.NewObjectID()
	}

	err := mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		_, err := ts.collection.InsertOne(sc, obj)
		return err
	})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (ts *transactionSessionImpl) UpdateByID(id primitive.ObjectID, update bson.M) (*bo.UsersBo, error) {
	ctx := context.Background()
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var result bo.UsersBo
	err := mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		return ts.collection.FindOneAndUpdate(sc, bson.M{"_id": id}, update, opts).Decode(&result)
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ts *transactionSessionImpl) DeleteByID(id primitive.ObjectID) error {
	ctx := context.Background()
	return mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		_, err := ts.collection.DeleteOne(sc, bson.M{"_id": id})
		return err
	})
}

func (ts *transactionSessionImpl) Count(filter bson.M) (int64, error) {
	ctx := context.Background()
	var count int64
	err := mongo.WithSession(ctx, ts.session, func(sc mongo.SessionContext) error {
		var err error
		count, err = ts.collection.CountDocuments(sc, filter)
		return err
	})
	return count, err
}
