package r34

import (
	"fmt"
	"r34-client/entities"
	"testing"
)

func TestR34Svc_GetPostByID(t *testing.T) {
	tc := []struct {
		ID     string
		Expect entities.PostDetail
	}{
		{ID: "9535856", Expect: entities.PostDetail{
			ID:          "9535856",
			Type:        entities.PostTypeImage,
			URL:         "https://wimg.rule34.xxx//samples/2052/sample_fc0e7cbb6d5c01e0a8164e810ed0219b.jpg?9535856",
			FullSizeURL: "https://wimg.rule34.xxx//images/2052/fc0e7cbb6d5c01e0a8164e810ed0219b.jpeg?9535856",
		}},
		{ID: "9578727", Expect: entities.PostDetail{
			ID:   "9578727",
			Type: entities.PostTypeVideo,
			URL:  "https://wimg.rule34.xxx//images/3584/c71b27944f4bdb96d352506e7af29ef5.jpg?9578727",
			// todo: sometime the domain is ahrimp4 and it can cause error, fix this
			FullSizeURL: "https://ahri2mp4.rule34.xxx//images/3584/c71b27944f4bdb96d352506e7af29ef5.mp4?9578727",
		}},
	}

	svc := New()
	for i, c := range tc {
		post, err := svc.GetPostByID(c.ID)
		if err != nil {
			t.Fatalf("Test case %d return error: %s", i, err)
		}

		if err := postDetailCompare(post, &c.Expect); err != nil {
			t.Logf("Test case %d FAILED, expected: %+v got %+v, err: %s", i, c.Expect, post, err)
			t.Fail()
		}

		t.Logf("Test case %d PASSED", i)
	}
}

func postDetailCompare(post *entities.PostDetail, expect *entities.PostDetail) error {
	if post.ID != expect.ID {
		return fmt.Errorf("ID not match")
	}
	if post.Type != expect.Type {
		return fmt.Errorf("Type not match")
	}
	if post.URL != expect.URL {
		return fmt.Errorf("URL not match")
	}
	if post.FullSizeURL != expect.FullSizeURL {
		return fmt.Errorf("FullSizeURL not match")
	}
	return nil
}

func TestR34Svc_ListPosts(t *testing.T) {
	tc := []struct {
		Param entities.GetPostsParams
		// im too lazy
		//ExpectedTotalPage uint
	}{
		{
			Param: entities.GetPostsParams{
				PaginationParam: entities.PaginationParam{},
				FilterParam:     entities.FilterParam{Search: "bigjohnson"},
			},
			//ExpectedTotalPage:
		},
	}

	svc := New()
	for i, s := range tc {
		posts, err := svc.ListPosts(&s.Param)
		if err != nil {
			t.Logf("Test case #%d return err: %s", i, err)
			t.Fail()
			continue
		}
		// todo: proper test
		t.Logf("Test case #%d output: %+v", i, posts)
		t.Logf("Test case #%d: PASS", i)
	}
}
