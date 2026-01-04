package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao/po"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UpdatePerson(c *gin.Context, in dto.UpdatePersonParaIn) (*dto.UpdatePersonParaOut, error) {
	// 0. 校验用户，需要鉴权jwt Token登录的必须验证正确性
	person, err := isCheckUserByGinContext(c)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// 只能修改自己的信息
	if person.ID != in.Id {
		err = fmt.Errorf("不能修改其他人的信息")
		zap.S().Error(err)
		return nil, err
	}

	// 1. 更新用户信息
	if in.Name != person.Name ||
		in.Idcard != person.Idcard ||
		in.Age != person.Age ||
		in.Gender != person.Gender ||
		in.Remark != person.Remark {
		tx := res.MysqlDB.Begin()
		fmt.Println("------update test_person表 事务开始------")

		personTxDao := dao.NewTestPersonDao(res.MysqlDB, tx, nil)
		param := bo.TestPersonBo{
			TestPerson: po.TestPerson{
				ID:     in.Id,
				Name:   in.Name,
				Idcard: in.Idcard,
				Age:    in.Age,
				Gender: in.Gender,
				Remark: in.Remark,
			},
		}
		personUpdate, err := personTxDao.Update(param)
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("更新失败")
			zap.S().Error("更新失败, 事务回滚, err: ", err)
			return nil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			err = fmt.Errorf("更新失败")
			zap.S().Error("更新失败, 事务回滚, err: ", err)
			return nil, err
		}
		fmt.Println("------update test_person表 事务提交------")

		return &dto.UpdatePersonParaOut{
			Person: dto.Person{
				Id:     personUpdate.ID,
				Name:   personUpdate.Name,
				Idcard: personUpdate.Idcard,
				Age:    personUpdate.Age,
				Gender: personUpdate.Gender,
			},
			Remark: personUpdate.Remark,
		}, nil
	} else {
		return &dto.UpdatePersonParaOut{
			Person: dto.Person{
				Id:     person.ID,
				Name:   person.Name,
				Idcard: person.Idcard,
				Age:    person.Age,
				Gender: person.Gender,
			},
			Remark: person.Remark,
		}, nil
	}
}
