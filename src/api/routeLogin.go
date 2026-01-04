package api

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/login"

	"github.com/gin-gonic/gin"
)

func routeLogin(r *gin.Engine) {
	gr1 := r.Group("/login")
	{
		gr1.POST("", login.Login)
	}

}
