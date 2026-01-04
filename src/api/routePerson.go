package api

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person"
	"github.com/go-framework-v2/go-backnormal-grpc/src/middleware"

	"github.com/gin-gonic/gin"
)

func routePerson(r *gin.Engine) {
	gr1 := r.Group("/person")
	{
		gr1.POST("/list", person.GetPersonList)
	}

	gr2 := r.Group("/person")
	{
		gr2.Use(middleware.JWTAuthMiddleware())

		gr2.POST("/insert", person.InsertPerson)
		gr2.POST("/update", person.UpdatePerson)
	}
}
