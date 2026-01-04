package dto

type UpdatePersonParaIn struct {
	Person
	Remark string `json:"remark"`
}

type UpdatePersonParaOut struct {
	Person
	Remark string `json:"remark"`
}
