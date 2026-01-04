package bo

import "github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/dao/po"

type UsersBo struct {
	po.Users `bson:",inline"` // 添加 inline 标签
}
