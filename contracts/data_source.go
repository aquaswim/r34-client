package contracts

import "r34-client/entities"

type DataSource interface {
	ListPosts(params *entities.GetPostsParams) (*entities.PostList, error)
	GetPostByID(id string) (*entities.PostDetail, error)
}
