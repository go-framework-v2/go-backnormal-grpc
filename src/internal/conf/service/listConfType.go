package service

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/dao"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/dto"
	"github.com/go-framework-v2/go-backnormal-grpc/src/res"

	"go.uber.org/zap"
)

func ListConfType(in dto.ListConfTypeParaIn) (*dto.ListConfTypeParaOut, error) {
	// default all
	confTypeDao := dao.NewConfTypeDao(res.MysqlDB, nil, nil)
	confTypeGetList, err := confTypeDao.GetAllConfTypeList()
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// convert to output
	outList := make([]dto.ConfType, 0)
	for _, confType := range confTypeGetList {
		outList = append(outList, dto.ConfType{
			TypeId:    confType.ID,
			TypeCode:  confType.TypeCode,
			TypeName:  confType.TypeName,
			Lifecycle: confType.Lifecycle,
			Category:  confType.Category,
			CreatedBy: confType.CreatedBy,
			UpdatedBy: confType.UpdatedBy,
			Remark:    confType.Remark,
			CreatedAt: confType.CreatedAt,
			UpdatedAt: confType.UpdatedAt,
			Valid:     confType.Valid,
		})
	}

	return &dto.ListConfTypeParaOut{
		List:  outList,
		Count: len(outList),
	}, nil
}
