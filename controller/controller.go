package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"path"
	"r34-client/contracts"
	"r34-client/entities"
	"r34-client/service/r34"
	lazyImage "r34-client/ui/lazy_image"
	"sync"
	"time"
)

type Controller struct {
	ListPostData    binding.UntypedList
	TotalPage       binding.Int
	CurrentPage     binding.Int
	dataSource      contracts.DataSource
	searchQuery     string
	listPostCache   sync.Map
	StatusText      binding.String
	DownloadPostIds binding.StringList
}

func New() *Controller {
	return &Controller{
		ListPostData:    binding.NewUntypedList(),
		TotalPage:       binding.NewInt(),
		CurrentPage:     binding.NewInt(),
		StatusText:      binding.NewString(),
		DownloadPostIds: binding.NewStringList(),
		dataSource:      r34.New(),
		listPostCache:   sync.Map{},
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

func (c *Controller) GetAutoComplete(q string) []string {
	c.SetStatusText("Getting auto complete for: " + q)
	resp, err := c.dataSource.GetAutoComplete(q)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("failed getting auto complete: %s", err))
	}
	return resp
}

func (c *Controller) ClearCache() {
	c.SetStatusText("Clearing all cache")
	// do clear all cache here
	err := lazyImage.ClearCache()
	if err != nil {
		c.SetStatusText(fmt.Sprintf("Error clearing cache: %s", err))
	}

	// reset all param
	c.ListPostData.Set([]interface{}{})
	c.TotalPage.Set(1)
	c.CurrentPage.Set(1)
	c.searchQuery = ""
	c.listPostCache = sync.Map{}
	c.SetStatusText("Cache cleared")
}

func (c *Controller) AddToDownloadList(postId string) {
	// get for duplicates
	existingIds, err := c.DownloadPostIds.Get()
	if err != nil {
		c.SetStatusText(fmt.Sprintf("internal error: %s", err))
	}
	for _, id := range existingIds {
		if id == postId {
			c.SetStatusText("post id already exists")
			return
		}
	}
	err = c.DownloadPostIds.Append(postId)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("internal error: %s", err))
		return
	}
	c.SetStatusText(fmt.Sprintf("post id %s added to download list", postId))
}

func (c *Controller) DownloadAllInList(outputFolder fyne.ListableURI) {
	if !c.IsDownloadListNotEmpty() {
		return
	}

	postIds, err := c.DownloadPostIds.Get()
	if err != nil {
		c.SetStatusText(fmt.Sprintf("internal error: %s", err))
	}

	details, err := c.getPostIdsDetails(postIds)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("Error getting post details: %s", err))
		return
	}

	for _, detail := range details {
		c.SetStatusText(fmt.Sprintf("processing postid: %s", detail.ID))
		srcUri, err := storage.ParseURI(detail.FullSizeURL)
		if err != nil {
			c.SetStatusText(fmt.Sprintf("Error parsing url of post %s, err: %s", detail.ID, err))
			return
		}
		outFile := path.Join(outputFolder.Path(), detail.ID+srcUri.Extension())
		writer, err := storage.Writer(storage.NewFileURI(outFile))
		if err != nil {
			c.SetStatusText(fmt.Sprintf("error write to file %s, err: %s", outFile, err))
		}
		c.downloadUri(srcUri, writer)
		// clear item from list
		c.removeDownloadListByValue(detail.ID)
	}
	c.SetStatusText("all item downloaded")
}

func (c *Controller) IsDownloadListNotEmpty() bool {
	return c.DownloadPostIds.Length() > 0
}

func (c *Controller) removeDownloadListByValue(id string) {
	postIds, err := c.DownloadPostIds.Get()
	if err != nil {
		c.SetStatusText(fmt.Sprintf("internal error: %s", err))
		return
	}
	newList := []string{}
	for _, postId := range postIds {
		if postId == id {
			continue
		}
		newList = append(newList, postId)
	}

	err = c.DownloadPostIds.Set(newList)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("internal error: %s", err))
		return
	}
}
