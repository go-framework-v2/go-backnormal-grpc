package api

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb"

	"github.com/gin-gonic/gin"
)

func routeMongodb(r *gin.Engine) {
	gr1 := r.Group("/mongodb")
	{
		gr1.POST("/users/list", mongodb.ListUsers)
	}

}
