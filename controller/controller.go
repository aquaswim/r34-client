package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"io"
	"r34-client/contracts"
	"r34-client/entities"
	"r34-client/service/r34"
	"sync"
	"time"
)

type Controller struct {
	ListPostData  binding.UntypedList
	TotalPage     binding.Int
	CurrentPage   binding.Int
	dataSource    contracts.DataSource
	searchQuery   string
	listPostCache sync.Map
	StatusText    binding.String
}

func New() *Controller {
	return &Controller{
		ListPostData:  binding.NewUntypedList(),
		TotalPage:     binding.NewInt(),
		CurrentPage:   binding.NewInt(),
		StatusText:    binding.NewString(),
		dataSource:    r34.New(),
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
		c.SetStatusText("list post received from cache")
		return val.(*entities.PostList), nil
	}

	c.SetStatusText("getting list post from server")
	d := time.Now()
	resp, err := c.dataSource.ListPosts(param)
	c.SetStatusText(fmt.Sprintf("list post received on: %s", time.Since(d)))

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
		c.SetStatusText(fmt.Sprintf("Error getting listpost with query %s: %s", c.searchQuery, err))
		return
	}
	listData := make([]interface{}, 0, response.TotalPage)
	for i, _ := range response.Items {
		listData = append(listData, &response.Items[i])
	}
	err = c.ListPostData.Set(listData)
	c.TotalPage.Set(int(response.TotalPage))
	if err != nil {
		c.SetStatusText(fmt.Sprintf("Error update list data binding: %s", err))
		return
	}
}

func (c *Controller) getCurrentPage() uint {
	p, _ := c.CurrentPage.Get()
	return uint(p)
}

func (c *Controller) SetStatusText(text string) {
	c.StatusText.Set(text)
}

func (c *Controller) OnClose() {

}

func (c *Controller) ResolvePostIDDownloadURI(id string) fyne.URI {
	c.SetStatusText(fmt.Sprintf("getting detail of postid: %s", id))
	postDetail, err := c.dataSource.GetPostByID(id)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error getting postId %s: %s", id, err))
		return nil
	}

	uri, err := storage.ParseURI(postDetail.FullSizeURL)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error parseurl %s: %s", postDetail.FullSizeURL, err))
		return nil
	}

	return uri
}

func (c *Controller) DownloadUri(srcUri fyne.URI, writer fyne.URIWriteCloser) {
	defer writer.Close()
	write, err := storage.CanWrite(writer.URI())
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error %s getting info writable for file: %s", writer.URI(), err))
		return
	}
	if write != true {
		c.SetStatusText(fmt.Sprintf("file %s not writeable", writer.URI()))
		return
	}

	c.SetStatusText(fmt.Sprintf("saving %s to %s", srcUri.Name(), writer.URI()))
	reader, err := storage.Reader(srcUri)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error reading %s: %s", reader.URI(), err))
		return
	}
	defer reader.Close()
	_, err = io.Copy(writer, reader)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error saving %s to %s", reader.URI(), writer.URI()))
		return
	}
	c.SetStatusText(fmt.Sprintf("%s saved to: %s", srcUri.Name(), writer.URI()))
}

func (c *Controller) GetAutoComplete(q string) []string {
	c.SetStatusText("Getting auto complete for: " + q)
	resp, err := c.dataSource.GetAutoComplete(q)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("failed getting auto complete: %s", err))
	}
	return resp
}
