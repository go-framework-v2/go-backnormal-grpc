package service

import (
	"context"
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListUsers(c *gin.Context, in dto.ListUsersParaIn) (*dto.ListUsersParaOut, error) {
	// 直接使用 MongoDB 客户端检查
	collection := res.MongoDB.Collection("users")

	// 查看集合中的所有文档
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Printf("Direct find error: %v\n", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rawDocs []bson.M
	if err := cursor.All(context.Background(), &rawDocs); err != nil {
		fmt.Printf("Cursor all error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Raw documents in collection: %d\n", len(rawDocs))
	for i, doc := range rawDocs {
		fmt.Printf("Raw doc[%d]: %+v\n", i, doc)
	}

	// 业务查询
	usersDao := dao.NewUsersDao(res.MongoDB)
	ids := []primitive.ObjectID{}
	idStrs := []string{"6901916c3340e2dc294f8801", "690191703340e2dc294f8802", "690191703340e2dc294f8803"}
	for _, idStr := range idStrs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			fmt.Printf("Invalid ObjectID: %s, error: %v\n", idStr, err)
			return nil, err
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		user, err := usersDao.FindByID(id)
		if err != nil {
			fmt.Printf("GetUserById error: %v\n", err)
			return nil, err
		}
		if user == nil {
			fmt.Printf("User not found with ID: %v\n", id)
			continue
		}
		fmt.Printf("User: %v\n", user)
	}

	return nil, nil
}
