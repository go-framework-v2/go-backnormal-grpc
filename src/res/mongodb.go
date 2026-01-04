package res

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-framework-v2/go-backnormal-grpc/src/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	// 根据你的实际模块路径修改
)

var (
	MongoDBClient *mongo.Client   // MongoDB 客户端实例
	MongoDB       *mongo.Database // MongoDB 数据库实例

	RESOURCE_MONGODB_NAME = "res-MongoDB" // res MongoDB数据库全局单例资源名称

	mongoDBOnce      sync.Once
	initMongoDBError error
)

// MongoDBResource 实现全局资源接口的包装器
type MongoDBResource struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (p *MongoDBResource) Name() string {
	return RESOURCE_MONGODB_NAME
}

func (p *MongoDBResource) Close() error {
	if p.Client != nil {
		return p.Client.Disconnect(context.Background())
	}
	return nil
}

// GetMongoDB 获取全局唯一的MongoDB数据库实例(线程安全)
func GetMongoDB() (*mongo.Database, error) {
	mongoDBOnce.Do(func() {
		client, db, err := initMongoDB()
		if err != nil {
			initMongoDBError = err
			return
		}
		MongoDBClient = client
		MongoDB = db
	})

	return MongoDB, initMongoDBError
}

// GetMongoDBClient 获取MongoDB客户端实例
func GetMongoDBClient() (*mongo.Client, error) {
	_, err := GetMongoDB()
	return MongoDBClient, err
}

// initMongoDB 初始化MongoDB连接
func initMongoDB() (*mongo.Client, *mongo.Database, error) {
	var err error

	host := config.Cfg.MongoDB.Host
	port := config.Cfg.MongoDB.Port
	db := config.Cfg.MongoDB.Db
	username := config.Cfg.MongoDB.Username
	password := config.Cfg.MongoDB.Password

	// 构建连接URI
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, db)

	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(dsn).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(5)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 建立连接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		zap.S().Errorf("Failed to connect to MongoDB: %v", err)
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		zap.S().Errorf("Failed to ping MongoDB: %v", err)
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// 选择数据库
	database := client.Database(db)

	// zap.S().Infof("SuccMongoDB resource registered successfullyessfully connected to MongoDB: %s/%s", host, db)
	return client, database, nil
}

// 新增：确保注册只执行一次
var registerMongoDBOnce sync.Once

// RegisterMongoDBToGlobal 将MongoDB实例注册到全局资源管理器
func RegisterMongoDBToGlobal() error {
	var registerErr error
	registerMongoDBOnce.Do(func() {
		// 先获取MongoDB实例，触发初始化
		client, err := GetMongoDBClient()
		if err != nil {
			registerErr = fmt.Errorf("failed to get MongoDB client: %w", err)
			zap.S().Error(registerErr)
			return
		}

		// 确保实例不为nil
		if client == nil || MongoDB == nil {
			registerErr = fmt.Errorf("MongoDB instances are nil")
			return
		}

		RegisterResource(RESOURCE_MONGODB_NAME, func() Resource {
			return &MongoDBResource{
				Client:   client,
				Database: MongoDB,
			}
		})

		// zap.L().Info("MongoDB resource registered successfully")
	})
	return registerErr
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return MongoDBClient.Disconnect(ctx)
	}
	return nil
}
