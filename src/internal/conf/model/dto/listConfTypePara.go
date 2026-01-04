package dto

type ListConfTypeParaIn struct{}

func (p ListConfTypeParaIn) Validate() error {
	return nil
}

type ListConfTypeParaOut struct {
	List  []ConfType `json:"list"`
	Count int        `json:"count"`
}
