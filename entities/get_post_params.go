package entities

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

func (p PaginationParam) Offset() uint {
	return p.PerPage * (p.Page - 1)
}
