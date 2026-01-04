package service

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/cons"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/dao/po"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/bo"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"go.uber.org/zap"
)

func AddConfType(in dto.AddConfTypeParaIn) (*dto.AddConfTypeParaOut, error) {
	typeCode := in.TypeCode
	typeName := in.TypeName
	category := in.Category
	createBy := in.CreatedBy
	updatedBy := in.UpdatedBy

	tx := res.MysqlDB.Begin()
	fmt.Println("------ AddConfType tx begin ------")
	confTypeDao := dao.NewConfTypeDao(res.MysqlDB, tx, nil)
	confTypeHisDao := dao.NewConfTypeHisDao(res.MysqlDB, tx, nil)

	// 1. add conf_type
	confTypeParam := bo.ConfTypeBo{
		ConfType: po.ConfType{
			TypeCode:  typeCode,
			TypeName:  typeName,
			Category:  category,
			CreatedBy: createBy,
			UpdatedBy: updatedBy,
		},
	}
	confTypeInsert, err := confTypeDao.Insert(confTypeParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("insert conf_type failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("insert confTypeInsert: ", confTypeInsert)

	// 2. add cong_type_his
	confTypeHisParam := bo.ConfTypeHisBo{
		ConfTypeHis: po.ConfTypeHis{
			TypeID:    confTypeInsert.ID,
			Operation: cons.CONF_TYPE_HIS_OPERATION_CREATE,
			NewValue:  fmt.Sprintf("%+v\n", confTypeInsert),
			CreatedBy: createBy,
			UpdatedBy: updatedBy,
		},
	}
	confTypeHisInsert, err := confTypeHisDao.Insert(confTypeHisParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("insert conf_type_his failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("insert confTypeHisInsert: ", confTypeHisInsert)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.S().Error("commit failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("------ AddConfType tx end ------")

	return &dto.AddConfTypeParaOut{
		ConfType: dto.ConfType{
			TypeId:    confTypeInsert.ID,
			TypeCode:  confTypeInsert.TypeCode,
			TypeName:  confTypeInsert.TypeName,
			Category:  confTypeInsert.Category,
			Lifecycle: confTypeInsert.Lifecycle,
			CreatedBy: confTypeInsert.CreatedBy,
			UpdatedBy: confTypeInsert.UpdatedBy,
			Remark:    confTypeInsert.Remark,
			CreatedAt: confTypeInsert.CreatedAt,
			UpdatedAt: confTypeInsert.UpdatedAt,
			Valid:     confTypeInsert.Valid,
		},
	}, nil
}
