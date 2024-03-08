package r34

import (
	"fmt"
	"net/url"
	"r34-client/entities"
	"strconv"
)

func getLastPageFromURL(urlStr string) (uint, error) {
	//https://rule34.xxx/index.php?page=post&s=list&tags=<not important>&pid=84 => 84 / max perpage
	parse, err := url.Parse(urlStr)
	if err != nil {
		return 0, err
	}
	pidStr := parse.Query().Get("pid")
	if pidStr == "" {
		return 0, fmt.Errorf("invalid last page url %s, pid query not detected", urlStr)
	}
	pidN, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}
	return uint((pidN / 42) + 1), nil
}

func generateListPostURL(params *entities.GetPostsParams) string {
	// https://rule34.xxx/index.php?page=post&s=list&tags=<search>&pid=<offset>
	q, _ := url.ParseQuery("")
	q.Set("page", "post")
	q.Set("s", "list")
	q.Set("tags", params.Search)
	q.Set("pid", strconv.Itoa(int(params.Offset())))

	return fmt.Sprintf("https://rule34.xxx/index.php?%s", q.Encode())
}
