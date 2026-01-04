package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/login/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao/po"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"
	dto2 "github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/service"
	"github.com/go-framework-v2/go-backnormal-grpc/src/middleware"

	"github.com/go-framework-v2/go-backnormal-grpc/src/tool"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(c *gin.Context, in dto.LoginParaIn) (*dto.LoginParaOut, error) {
	// name, idcard
	name := in.Name
	idcard := in.Idcard
	fmt.Printf("in.Name: %s, in.Idcard: %s\n", name, idcard)

	// ip
	ip := tool.GetIpByGinContext(c)
	fmt.Printf("ip: %s\n", ip)

	// 1. 根据name+idcard查询用户信息，不存在则插入
	param := bo.TestPersonBo{
		TestPerson: po.TestPerson{
			Name:   name,
			Idcard: idcard,
			Age:    in.Age,
			Gender: in.Gender,
		},
	}
	person, err := service.GetInsertPerson(param)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	fmt.Println("查询用户person: ", person)

	// 2. 生成token
	token, err := middleware.GenerateToken(person.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	fmt.Println("根据userId生成token: ", token)

	return &dto.LoginParaOut{
		Person: dto2.Person{
			Id:     person.ID,
			Name:   person.Name,
			Idcard: person.Idcard,
			Age:    person.Age,
			Gender: person.Gender,
		},
		Token: token,
	}, nil
}
