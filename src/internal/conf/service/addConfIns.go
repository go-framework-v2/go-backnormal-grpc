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

func AddConfIns(in dto.AddConfInsParaIn) (*dto.AddConfInsParaOut, error) {
	typeId := in.TypeId
	typeCode := in.TypeCode
	typeName := in.TypeName
	category := in.Category
	// 1. validate conf_type
	confTypeDao := dao.NewConfTypeDao(res.MysqlDB, nil, nil)
	confTypeParam := bo.ConfTypeBo{
		ConfType: po.ConfType{
			ID:       typeId,
			TypeCode: typeCode,
			TypeName: typeName,
			Category: category,
		},
	}
	confTypeGet, err := confTypeDao.GetConfType(confTypeParam)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	if confTypeGet == nil {
		err = fmt.Errorf("confType not found")
		zap.S().Error(err)
		return nil, err
	}

	insCode := in.InsCode
	insName := in.InsName
	remark := in.Remark
	createBy := in.CreatedBy
	updatedBy := in.UpdatedBy

	tx := res.MysqlDB.Begin()
	fmt.Println("------ AddConfIns tx begin ------")
	confInsDao := dao.NewConfInsDao(res.MysqlDB, tx, nil)
	confInsHisDao := dao.NewConfInsHisDao(res.MysqlDB, tx, nil)

	parentInsId := in.ParentInsId
	parentInsCode := in.ParentInsCode
	parentInsName := in.ParentInsName
	if parentInsId > 0 && parentInsCode != "" && parentInsName != "" {
		insCode = fmt.Sprintf("%s.%s", parentInsCode, insCode) // 拼接, 子实例编码格式: 父实例编码.子实例编码
		fmt.Println("insCode: ", insCode)
	}
	// 2. insert conf_ins
	confInsParam := bo.ConfInsBo{
		ConfIns: po.ConfIns{
			InsCode:   insCode,
			InsName:   insName,
			TypeID:    typeId,
			CreatedBy: createBy,
			UpdatedBy: updatedBy,
			Remark:    remark,
		},
	}
	confInsInsert, err := confInsDao.Insert(confInsParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("insert conf_ins failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("insert confInsInsert: ", confInsInsert)

	// 3. insert conf_ins_his
	confInsHisParam := bo.ConfInsHisBo{
		ConfInsHis: po.ConfInsHis{
			InsID:     confInsInsert.ID,
			Operation: cons.CONF_TYPE_HIS_OPERATION_CREATE,
			NewValue:  fmt.Sprintf("%+v\n", confInsInsert),
			CreatedBy: createBy,
			UpdatedBy: updatedBy,
		},
	}
	confInsHisInsert, err := confInsHisDao.Insert(confInsHisParam)
	if err != nil {
		tx.Rollback()
		zap.S().Error("insert conf_ins_his failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("insert confInsHisInsert: ", confInsHisInsert)

	if parentInsId > 0 && parentInsCode != "" && parentInsName != "" {
		// 4. insert conf
		// 4.1 get parent conf_ins
		parentConfInsParam := bo.ConfInsBo{
			ConfIns: po.ConfIns{
				ID:      parentInsId,
				InsCode: parentInsCode,
				InsName: parentInsName,
			},
		}
		parentConfInsGet, err := confInsDao.GetConfIns(parentConfInsParam)
		if err != nil {
			tx.Rollback()
			zap.S().Error("get parent conf_ins failed, tx rollback, err: ", err)
			return nil, err
		}
		if parentConfInsGet == nil {
			tx.Rollback()
			err = fmt.Errorf("parent conf_ins not found")
			zap.S().Error(err)
			return nil, err
		}

		confDao := dao.NewConfDao(res.MysqlDB, tx, nil)
		confHisDao := dao.NewConfHisDao(res.MysqlDB, tx, nil)
		// 4.2 insert conf
		confParam := bo.ConfBo{
			Conf: po.Conf{
				SourceInsID: parentInsId,
				TargetInsID: confInsInsert.ID,
				RefType:     cons.CONF_REF_TYPE_CONTAINS,
				CreatedBy:   createBy,
				UpdatedBy:   updatedBy,
			},
		}
		confInsert, err := confDao.Insert(confParam)
		if err != nil {
			tx.Rollback()
			zap.S().Error("insert conf failed, tx rollback, err: ", err)
			return nil, err
		}
		fmt.Println("insert confInsert: ", confInsert)

		// 4.3 insert conf_his
		confHisParam := bo.ConfHisBo{
			ConfHis: po.ConfHis{
				ConfID:    confInsert.ID,
				Operation: cons.CONF_HIS_OPERATION_CREATE,
				NewValue:  fmt.Sprintf("%+v\n", confInsert),
				CreatedBy: createBy,
				UpdatedBy: updatedBy,
			},
		}
		confHisInsert, err := confHisDao.Insert(confHisParam)
		if err != nil {
			tx.Rollback()
			zap.S().Error("insert conf_his failed, tx rollback, err: ", err)
			return nil, err
		}
		fmt.Println("insert confHisInsert: ", confHisInsert)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		zap.S().Error("commit failed, tx rollback, err: ", err)
		return nil, err
	}
	fmt.Println("------ AddConfIns tx end ------")

	return &dto.AddConfInsParaOut{
		ConfIns: dto.ConfIns{
			InsId:     confInsInsert.ID,
			InsCode:   confInsInsert.InsCode,
			InsName:   confInsInsert.InsName,
			TypeId:    confTypeGet.ID,
			CreatedBy: confInsInsert.CreatedBy,
			UpdatedBy: confInsInsert.UpdatedBy,
			Remark:    confInsInsert.Remark,
			CreatedAt: confInsInsert.CreatedAt,
			UpdatedAt: confInsInsert.UpdatedAt,
			Valid:     confInsInsert.Valid,
		},
	}, nil
}
