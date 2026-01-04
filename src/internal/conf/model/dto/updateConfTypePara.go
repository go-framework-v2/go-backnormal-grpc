package dto

import "fmt"

type UpdateConfTypeParaIn struct {
	// 必填
	TypeId int32 `json:"type_id"`
	// 可选
	TypeCode string `json:"type_code"`
	TypeName string `json:"type_name"`
	Category string `json:"category"`
	Remark   string `json:"remark"`
}

func (p UpdateConfTypeParaIn) Validate() error {
	if p.TypeId <= 0 {
		return fmt.Errorf("type_id must be greater than 0")
	}

	return nil
}

type UpdateConfTypeParaOut struct {
	ConfType ConfType
}
