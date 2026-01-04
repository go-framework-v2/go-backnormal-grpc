package po

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionNameUsers = "users"

// Users MongoDB 文档结构
type Users struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Age int32 `bson:"age"`
	Email string `bson:"email"`
	Name string `bson:"name"`
}

// CollectionName 返回集合名称
func (*Users) CollectionName() string {
	return CollectionNameUsers
}
