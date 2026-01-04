package mongodb

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	tool.HandleWithBindWithC(
		c,
		service.ListUsers,
		dto.ListUsersParaOut{},
	)
}
