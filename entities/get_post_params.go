package entities

import "encoding/json"

type PaginationParam struct {
	PerPage uint
	Page    uint
}

type FilterParam struct {
	Search string
}

type GetPostsParams struct {
	PaginationParam
	FilterParam
}

func (p GetPostsParams) ToHash() (string, error) {
	marshaled, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}
	return string(marshaled), nil
}

func (p PaginationParam) Offset() uint {
	return p.PerPage * (p.Page - 1)
}
