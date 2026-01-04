package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"go.uber.org/zap"
)

func GetPersonList(in dto.GetPersonListParaIn) (*dto.GetPersonListParaOut, error) {
	page := in.Page
	pageSize := in.PageSize

	testPersonDao := dao.NewTestPersonDao(res.MysqlDB, nil, nil)
	param := bo.TestPersonBo{
		Page:     page,
		PageSize: pageSize,
	}
	testPersonList, err := testPersonDao.GetTestPersonList(param)
	if err != nil {
		err = fmt.Errorf("查询错误")
		zap.S().Errorf("查询错误, err: %v", err)
		return nil, err
	}

	var list []dto.Person
	for _, v := range testPersonList {
		list = append(list, dto.Person{
			Id:     v.ID,
			Name:   v.Name,
			Idcard: v.Idcard,
			Age:    v.Age,
			Gender: v.Gender,
		})
	}
	return &dto.GetPersonListParaOut{
		List:  list,
		Count: len(list),
	}, nil
}
