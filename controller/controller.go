package controller

import (
	"fyne.io/fyne/v2/data/binding"
	"log"
	"r34-client/commons"
	"r34-client/contracts"
	"r34-client/entities"
	"r34-client/service/r34"
	"sync"
)

type Controller struct {
	ListPostData  binding.UntypedList
	TotalPage     binding.Int
	CurrentPage   binding.Int
	dataSource    contracts.DataSource
	searchQuery   string
	l             *log.Logger
	listPostCache sync.Map
}

func New() *Controller {
	return &Controller{
		ListPostData:  binding.NewUntypedList(),
		TotalPage:     binding.NewInt(),
		CurrentPage:   binding.NewInt(),
		dataSource:    r34.New(),
		l:             commons.NewLogger("[CTRL] "),
		listPostCache: sync.Map{},
	}
}

func (c *Controller) Search(query string) {
	// reset controller property
	c.searchQuery = query
	c.TotalPage.Set(1)
	c.CurrentPage.Set(1)

	c.fetchPosts()
}

func (c *Controller) ChangePage(page int) {
	c.CurrentPage.Set(page)
	c.fetchPosts()
}

func (c *Controller) getListPostFromDataSource(param *entities.GetPostsParams) (*entities.PostList, error) {
	paramHash, err := param.ToHash()
	if err != nil {
		return nil, err
	}
	// check if it exist in cache
	val, ok := c.listPostCache.Load(paramHash)
	if ok {
		return val.(*entities.PostList), nil
	}

	resp, err := c.dataSource.ListPosts(param)

	// save to cache when not error
	if err == nil {
		c.listPostCache.Store(paramHash, resp)
	}

	return resp, err
}

func (c *Controller) fetchPosts() {
	c.ListPostData.Set(make([]interface{}, 0))
	response, err := c.getListPostFromDataSource(&entities.GetPostsParams{
		PaginationParam: entities.PaginationParam{
			PerPage: 42,
			Page:    c.getCurrentPage(),
		},
		FilterParam: entities.FilterParam{Search: c.searchQuery},
	})
	if err != nil {
		log.Printf("Error getting listpost with query %s: %s", c.searchQuery, err)
		return
	}
	listData := make([]interface{}, 0, response.TotalPage)
	for i, _ := range response.Items {
		listData = append(listData, &response.Items[i])
	}
	err = c.ListPostData.Set(listData)
	c.TotalPage.Set(int(response.TotalPage))
	if err != nil {
		log.Printf("Error update list data binding: %s", err)
		return
	}
}

func (c *Controller) getCurrentPage() uint {
	p, _ := c.CurrentPage.Get()
	return uint(p)
}

func (c *Controller) OnClose() {

}
