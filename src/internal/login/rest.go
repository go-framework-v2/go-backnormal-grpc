package login

import (
	"net/http"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/login/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/login/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
	"github.com/go-framework-v2/go-access/access"
	"go.uber.org/zap"
)

func Login(c *gin.Context) {
	zap.S().Debug("------------用户登录接口------------")
	// 错误处理，返回各字段的零值，避前端处理报错
	emptyData := dto.LoginParaOut{}

	// 1. 处理请求体(JSON格式) + 验证
	in, err := tool.ShouldBindJSON_Validate[dto.LoginParaIn](c)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 2. 处理业务逻辑
	outData, err := service.Login(c, in.Data)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 3. 返回响应数据(JSON格式)
	out := access.GetSuccessResult(outData, "Success")
	c.JSON(http.StatusOK, out)
}
