package dto

import (
	"fmt"

	"github.com/go-framework-v2/go-access/common"
)

type GetPersonListParaIn struct {
	common.Paginator
}

func (p *GetPersonListParaIn) Validate() error {
	if p.Page <= 0 {
		return fmt.Errorf("page must be greater than 0")
	}
	if p.PageSize <= 0 {
		return fmt.Errorf("pageSize must be greater than 0")
	}
	return nil
}

type GetPersonListParaOut struct {
	List  []Person `json:"list"`
	Count int      `json:"count"`
}

type Person struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Idcard string `json:"idcard"`
	Age    int32  `json:"age"`
	Gender string `json:"gender"`
}
