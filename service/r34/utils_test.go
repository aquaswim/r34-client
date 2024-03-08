package r34

import "testing"

func TestUtils_getLastPageFromURL(t *testing.T) {
	tc := []struct {
		Url      string
		Expected uint
	}{
		{Url: "https://rule34.xxx/index.php?page=post&s=list&tags=test&pid=84", Expected: 3},
		{Url: "https://rule34.xxx/index.php?page=post&s=list&tags=test+test2&pid=42", Expected: 2},
	}

	for i, s := range tc {
		output, err := getLastPageFromURL(s.Url)
		if err != nil {
			t.Fatalf("Test case #%d: getLastPageFromURL return err: %s", i, err)
		}
		if output != s.Expected {
			t.Logf("Test case #%d: FAIL expected %d got %d", i, s.Expected, output)
			t.Fail()
		} else {
			t.Logf("Test case #%d: PASS", i)
		}
	}
}
