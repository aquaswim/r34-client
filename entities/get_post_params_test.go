package entities

import "testing"

func TestPaginationParam_Offset(t *testing.T) {
	tc := []struct {
		Param          PaginationParam
		ExpectedOffset uint
	}{
		{Param: PaginationParam{
			PerPage: 0,
			Page:    0,
		}, ExpectedOffset: 0},
		{Param: PaginationParam{
			PerPage: 10,
			Page:    1,
		}, ExpectedOffset: 0},
		{Param: PaginationParam{
			PerPage: 10,
			Page:    2,
		}, ExpectedOffset: 10},
		{Param: PaginationParam{
			PerPage: 42,
			Page:    2,
		}, ExpectedOffset: 42},
	}

	for i, s := range tc {
		if s.Param.Offset() != s.ExpectedOffset {
			t.Logf("Test case #%d failed, %+v expected offset is %d but got %d", i, s.Param, s.ExpectedOffset, s.Param.Offset())
		}
		t.Logf("Test case #%d PASS", i)
	}
}
