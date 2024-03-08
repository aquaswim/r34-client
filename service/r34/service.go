package r34

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"r34-client/commons"
	"r34-client/contracts"
	"r34-client/entities"
	"strconv"
)

type r34Svc struct {
	l *log.Logger
}

func New() contracts.DataSource {
	return &r34Svc{
		l: commons.NewLogger("r34Svc"),
	}
}

func (r r34Svc) ListPosts(params *entities.GetPostsParams) (*entities.PostList, error) {
	// maximum perpage in r34 is 42
	if params.PerPage > 42 {
		return nil, fmt.Errorf("maximum perpage is 42")
	}

	// https://rule34.xxx/index.php?page=post&s=list&tags=<search>&pid=<offset>
	//url := fmt.Sprintf("https://rule34.xxx/index.php?page=post&s=list&tags=%s&pid=%d", params.Search, params.Offset())
	url := generateListPostURL(params)
	r.l.Printf("get url %s", url)

	doc, err := r.getDocumentFromURL(url)
	if err != nil {
		r.l.Printf("error parsing response %s, error: %s", url, err)
		return nil, err
	}

	// getting post count
	postsElem := doc.Find(".image-list span.thumb")
	totalPostPerPage := postsElem.Length()

	if totalPostPerPage < 1 {
		err = fmt.Errorf("total post perpage in %s is less than zero", url)
		r.l.Println(err)
		return nil, err
	}
	// getting pagination total page
	// last paginator is either anchor with ">>" text or "b" if already on last page
	lastPaginator := doc.Find("#paginator .pagination a,b").Last()
	var totalPage uint = 1
	if lastpageUrl := lastPaginator.AttrOr("href", ""); lastPaginator.Text() == ">>" && lastpageUrl != "" {
		totalPageFromUrl, err := getLastPageFromURL(lastPaginator.AttrOr("href", ""))
		if err != nil {
			r.l.Printf("error getting last page from url: %s, error: %s", lastpageUrl, err)
		} else {
			totalPage = totalPageFromUrl
		}
	} else if lastPageInt, err := strconv.Atoi(lastPaginator.Text()); err == nil {
		totalPage = uint(lastPageInt)
	}

	var list = make([]entities.Post, totalPostPerPage)
	postsElem.Each(func(i int, el *goquery.Selection) {
		// id always start with s, delete that from string
		list[i].ID = el.AttrOr("id", "s")[1:]
		imgEle := el.Find("img")
		list[i].ThumbnailURL = imgEle.AttrOr("src", "")
		// video post always have this style attr: border: 3px solid rgb(0, 0, 255);
		if imgEle.AttrOr("style", "") == "" {
			list[i].Type = entities.PostTypeImage
		} else {
			list[i].Type = entities.PostTypeVideo
		}
	})

	return &entities.PostList{
		Items:     list,
		TotalPage: totalPage,
	}, nil
}

func (r r34Svc) GetPostByID(id string) (*entities.PostDetail, error) {
	// https://rule34.xxx/index.php?page=post&s=view&id=<id>
	url := "https://rule34.xxx/index.php?page=post&s=view&id=" + id
	r.l.Printf("get url %s", url)

	doc, err := r.getDocumentFromURL(url)
	if err != nil {
		r.l.Printf("error parsing response %s, error: %s", url, err)
		return nil, err
	}
	parsed := new(entities.PostDetail)
	// check is content type
	//  image = have #image
	//  video = tba
	if srcUrl, found := doc.Find("#image").Attr("src"); found {
		parsed.Type = entities.PostTypeImage
		parsed.ID, _ = doc.Find("#edit_form input[name=id]").Attr("value")
		parsed.URL = srcUrl
		parsed.FullSizeURL, _ = doc.Find("meta[itemprop=image]").Attr("content")
	} else if srcUrl, found := doc.Find("video").Attr("poster"); found {
		parsed.Type = entities.PostTypeVideo
		parsed.ID, _ = doc.Find("#edit_form input[name=id]").Attr("value")
		parsed.URL = srcUrl
		parsed.FullSizeURL, _ = doc.Find("video source").Attr("src")
	}

	return parsed, nil
}

func (r r34Svc) getDocumentFromURL(urlStr string) (*goquery.Document, error) {
	res, err := http.Get(urlStr)
	if err != nil {
		r.l.Printf("error getting from url: %s, error: %s", urlStr, err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("getting url %s return status code error: %d %s", urlStr, res.StatusCode, res.Status)
		r.l.Println(err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		r.l.Printf("error parsing response %s, error: %s", urlStr, err)
		return nil, err
	}
	return doc, nil
}
