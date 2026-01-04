package dto

import (
	"fmt"

	"github.com/go-framework-v2/go-backnormal-grpc/src/cons"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/dto"
)

type LoginParaIn struct {
	Name   string `json:"name"`
	Idcard string `json:"idcard"`
	Age    int32  `json:"age"`
	Gender string `json:"gender"`
}

func (p LoginParaIn) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Idcard == "" {
		return fmt.Errorf("idcard is required")
	}

	if p.Age < 0 {
		return fmt.Errorf("age is invalid")
	}
	if p.Gender != cons.GENDER_MALE && p.Gender != cons.GENDER_FEMALE {
		return fmt.Errorf("gender is invalid")
	}

	return nil
}

type LoginParaOut struct {
	dto.Person
	Token string `json:"token"`
}
