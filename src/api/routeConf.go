package api

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf"

	"github.com/gin-gonic/gin"
)

func routeConf(r *gin.Engine) {
	gr1 := r.Group("/confType")
	{
		gr1.POST("/add", conf.AddConfType)       // 插入
		gr1.POST("/update", conf.UpdateConfType) // 更新
		gr1.POST("/delete", conf.DeleteConfType) // 删除(逻辑)
		gr1.POST("/list", conf.ListConfType)     // 查询
	}

	gr2 := r.Group("/confIns")
	{
		gr2.POST("/add", conf.AddConfIns)       // 插入
		gr2.POST("/update", conf.UpdateConfIns) // 更新
		gr2.POST("/delete", conf.DeleteConfIns) // 删除(逻辑)
		gr2.POST("/list", conf.ListConfIns)     // 查询
	}
}
