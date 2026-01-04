package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/cons"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"go.uber.org/zap"
)

// GetInsertPerson 根据用户信息(主键信息)获取用户信息，如果不存在，则插入用户信息
func GetInsertPerson(param bo.TestPersonBo) (*bo.TestPersonBo, error) {
	// 1. 根据主键信息查询用户信息
	name := param.Name
	idcard := param.Idcard
	gender := param.Gender
	if name == "" || idcard == "" {
		err := fmt.Errorf("参数错误")
		zap.S().Error("name or idcard is empty")
		return nil, err
	}
	if gender != cons.GENDER_MALE && gender != cons.GENDER_FEMALE {
		err := fmt.Errorf("参数错误")
		zap.S().Error("gender is not valid")
		return nil, err
	}

	testPersonDao := dao.NewTestPersonDao(res.MysqlDB, nil, nil)
	person, err := testPersonDao.FindOneByUk(name, idcard)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// 2. 如果用户信息存在，则直接返回用户信息
	if person != nil {
		return person, nil
	}

	// 3. 如果用户信息不存在，则插入用户信息
	tx := res.MysqlDB.Begin()
	fmt.Println("------insert test_person表 开始事务------")

	testPersonTxDao := dao.NewTestPersonDao(res.MysqlDB, tx, nil)
	person, err = testPersonTxDao.Insert(param)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("插入失败")
		zap.S().Error("插入失败, 事务回滚, err: ", err)
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		err = fmt.Errorf("提交事务失败")
		zap.S().Error("提交事务失败, err: ", err)
		return nil, err
	}
	fmt.Println("------insert test_person表 事务提交成功------")

	// 4. 查询事务提交之后的用户信息
	id := person.ID
	person, err = testPersonDao.FindOne(int(id))
	if err != nil {
		err = fmt.Errorf("查询失败")
		zap.S().Error("查询失败, err: ", err)
		return nil, err
	}

	return person, nil
}
