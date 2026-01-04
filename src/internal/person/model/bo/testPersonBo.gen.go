package bo

import "github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/dao/po"

type TestPersonBo struct {
	po.TestPerson
	Page     int `gorm:"-"`
	PageSize int `gorm:"-"`
}
