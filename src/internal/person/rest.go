package person

import (
	"net/http"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
	"github.com/go-framework-v2/go-access/access"
)

func GetPersonList(c *gin.Context) {
	// 错误处理，返回各字段的零值，避前端处理报错
	emptyData := dto.GetPersonListParaOut{
		List:  []dto.Person{}, // []而不是null
		Count: 0,
	}

	// 1. 处理请求体(JSON格式) + 验证
	in, err := tool.ShouldBindJSON_Validate[dto.GetPersonListParaIn](c)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 2. 处理业务逻辑
	outData, err := service.GetPersonList(in.Data)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 3. 返回响应数据(JSON格式)
	out := access.GetSuccessResult(outData, "Success")
	c.JSON(http.StatusOK, out)
}

func InsertPerson(c *gin.Context) {

}

func UpdatePerson(c *gin.Context) {
	// 错误处理，返回各字段的零值，避前端处理报错
	emptyData := dto.UpdatePersonParaIn{}

	// 1. 处理请求体(JSON格式) + 验证
	in, err := tool.ShouldBindJSON_Validate[dto.UpdatePersonParaIn](c)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 2. 处理业务逻辑
	outData, err := service.UpdatePerson(c, in.Data)
	if err != nil {
		out := access.GetSuccessResult(emptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 3. 返回响应数据(JSON格式)
	out := access.GetSuccessResult(outData, "Success")
	c.JSON(http.StatusOK, out)
}
