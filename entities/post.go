package entities

type Post struct {
	ID           string
	ThumbnailURL string
	Type         PostType
}

type PostType int8

const (
	PostTypeImage PostType = 1
	PostTypeVideo PostType = 2
)

type PostDetail struct {
	ID          string
	Type        PostType
	URL         string
	FullSizeURL string
}

type PostList struct {
	Items     []Post
	TotalPage uint
}
