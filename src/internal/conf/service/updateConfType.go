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

func UpdateConfType(in dto.UpdateConfTypeParaIn) (*dto.UpdateConfTypeParaOut, error) {
	typeId := in.TypeId
	typeCode := in.TypeCode
	typeName := in.TypeName
	category := in.Category
	remark := in.Remark

	// 0. get conf_type
	confTypeNoTxDao := dao.NewConfTypeDao(res.MysqlDB, nil, nil)
	confTypeGet, err := confTypeNoTxDao.FindOne(int(typeId))
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	if confTypeGet == nil {
		err = fmt.Errorf("conf_type not found, type_id: %d", typeId)
		zap.S().Error(err)
		return nil, err
	}
	createBy := confTypeGet.CreatedBy
	updatedBy := confTypeGet.UpdatedBy

	tx := res.MysqlDB.Begin()
	fmt.Println("------ UpdateConfType tx begin ------")
	confTypeDao := dao.NewConfTypeDao(res.MysqlDB, tx, nil)
	confTypeHisDao := dao.NewConfTypeHisDao(res.MysqlDB, tx, nil)

	// 1. update conf_type
	confTypeParam := bo.ConfTypeBo{
		ConfType: po.ConfType{
			ID:       typeId,
			TypeCode: typeCode,
			TypeName: typeName,
			Category: category,
			Remark:   remark,
		},
	}
	confTypeUpdate, err := confTypeDao.Update(confTypeParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("update conf_type failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("update confTypeUpdate: ", confTypeUpdate)

	// 2. insert cong_type_his
	confTypeHisParam := bo.ConfTypeHisBo{
		ConfTypeHis: po.ConfTypeHis{
			TypeID:    typeId,
			Operation: cons.CONF_TYPE_HIS_OPERATION_UPDATE,
			OldValue:  fmt.Sprintf("%+v\n", confTypeGet),
			NewValue:  fmt.Sprintf("%+v\n", confTypeUpdate),
			CreatedBy: createBy,
			UpdatedBy: updatedBy,
		},
	}
	confTypeHisUpdate, err := confTypeHisDao.Insert(confTypeHisParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("insert conf_type_his failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("insert confTypeHisInsert: ", confTypeHisUpdate)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.S().Error("commit failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("------ UpdateConfType tx end ------")

	return &dto.UpdateConfTypeParaOut{
		ConfType: dto.ConfType{
			TypeId:    confTypeUpdate.ID,
			TypeCode:  confTypeUpdate.TypeCode,
			TypeName:  confTypeUpdate.TypeName,
			Category:  confTypeUpdate.Category,
			Lifecycle: confTypeUpdate.Lifecycle,
			CreatedBy: confTypeUpdate.CreatedBy,
			UpdatedBy: confTypeUpdate.UpdatedBy,
			Remark:    confTypeUpdate.Remark,
			CreatedAt: confTypeUpdate.CreatedAt,
			UpdatedAt: confTypeUpdate.UpdatedAt,
			Valid:     confTypeUpdate.Valid,
		},
	}, nil
}
