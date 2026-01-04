package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// isCheckUserByGinContext 根据gin.Context获取用户信息，校验是否存在且为stepcash应用用户
func isCheckUserByGinContext(c *gin.Context) (*bo.TestPersonBo, error) {
	// 从鉴权中间件获取userId
	cUserId, exists := c.Get("userId")
	if !exists {
		err := fmt.Errorf("server has no userId")
		zap.L().Error(err.Error())
		return nil, err
	}
	userId := cUserId.(int64) // unit64默认设置的数据类型

	// 校验用户是否存在
	personDao := dao.NewTestPersonDao(res.MysqlDB, nil, nil)
	person, err := personDao.FindOne(int(userId))
	if err != nil {
		err := fmt.Errorf("查询失败")
		zap.S().Error("查询失败, err: ", err)
		return nil, err
	}
	if person == nil {
		err := fmt.Errorf("用户不存在")
		zap.S().Error("用户不存在, err: ", err)
		return nil, err
	}

	return person, nil
}
