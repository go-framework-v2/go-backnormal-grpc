package dto

import (
	"fmt"
	"time"
)

type AddConfTypeParaIn struct {
	// 必填
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
	Category string `json:"category"`
	// 可选
	CreatedBy string `json:"create_by"`
	UpdatedBy string `json:"updated_by"`
}

func (p AddConfTypeParaIn) Validate() error {
	if p.TypeCode == "" {
		return fmt.Errorf("TypeCode is required")
	}
	if p.TypeName == "" {
		return fmt.Errorf("TypeName is required")
	}
	if p.Category == "" {
		return fmt.Errorf("Category is required")
	}

	return nil
}

type AddConfTypeParaOut struct {
	ConfType ConfType
}

type ConfType struct {
	TypeId    int32     `json:"type_id"`
	TypeCode  string    `json:"type_code"`
	TypeName  string    `json:"type_name"`
	Lifecycle string    `json:"lifecycle"`
	Category  string    `json:"category"`
	CreatedBy string    `json:"create_by"`
	UpdatedBy string    `json:"updated_by"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Valid     int32     `json:"valid"`
}
