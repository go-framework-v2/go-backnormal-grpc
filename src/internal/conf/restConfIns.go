package conf

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
)

func AddConfIns(c *gin.Context) {
	tool.HandleWithBind(
		c,
		service.AddConfIns,
		dto.AddConfInsParaOut{},
	)
}

func UpdateConfIns(c *gin.Context) {
	tool.HandleWithBind(
		c,
		service.UpdateConfIns,
		dto.UpdateConfInsParaOut{},
	)
}

func DeleteConfIns(c *gin.Context) {

}

func ListConfIns(c *gin.Context) {
	tool.HandleWithBind(
		c,
		service.ListConfIns,
		dto.ListConfInsParaOut{},
	)
}
