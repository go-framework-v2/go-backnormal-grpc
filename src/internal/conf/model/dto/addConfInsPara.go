package dto

import (
	"fmt"
	"time"
)

type AddConfInsParaIn struct {
	// 必填
	// 配置实体分类
	TypeId   int32  `json:"type_id"`
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
	Category string `json:"category"`
	// 配置实体
	InsCode string `json:"ins_code"`
	InsName string `json:"ins_name"`

	// 可选
	Remark string `json:"remark"`
	// 配置实体创建人
	CreatedBy string `json:"create_by"`
	UpdatedBy string `json:"updated_by"`
	// 配置实体父实体
	ParentInsId   int64  `json:"parent_ins_id"`
	ParentInsCode string `json:"parent_ins_code"`
	ParentInsName string `json:"parent_ins_name"`
}

func (p AddConfInsParaIn) Validate() error {
	if p.TypeId <= 0 {
		return fmt.Errorf("type_id must be greater than 0")
	}

	if p.InsCode == "" {
		return fmt.Errorf("ins_code must not be empty")
	}
	if p.InsName == "" {
		return fmt.Errorf("ins_name must not be empty")
	}

	return nil
}

type AddConfInsParaOut struct {
	ConfIns ConfIns
}

type ConfIns struct {
	InsId     int64     `json:"ins_id"`
	InsCode   string    `json:"ins_code"`
	InsName   string    `json:"ins_name"`
	TypeId    int32     `json:"type_id"`
	CreatedBy string    `json:"create_by"`
	UpdatedBy string    `json:"updated_by"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Valid     int32     `json:"valid"`
}
