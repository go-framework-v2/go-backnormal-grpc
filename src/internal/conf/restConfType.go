package conf

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
)

func AddConfType(c *gin.Context) {
	tool.HandleWithBind[dto.AddConfTypeParaIn, dto.AddConfTypeParaOut](
		c,
		service.AddConfType,
		dto.AddConfTypeParaOut{},
	)
}

func UpdateConfType(c *gin.Context) {
	tool.HandleWithBind[dto.UpdateConfTypeParaIn, dto.UpdateConfTypeParaOut](
		c,
		service.UpdateConfType,
		dto.UpdateConfTypeParaOut{},
	)
}

func DeleteConfType(c *gin.Context) {

}

func ListConfType(c *gin.Context) {
	tool.HandleWithBind[dto.ListConfTypeParaIn, dto.ListConfTypeParaOut](
		c,
		service.ListConfType,
		dto.ListConfTypeParaOut{},
	)
}
