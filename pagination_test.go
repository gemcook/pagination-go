package pager

import (
	"reflect"
	"testing"
)

func TestParsePagination(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want *Pagination
	}{
		// TODO: Add test cases.
		{"default limit=30, page=1", args{"https://example.com/fruits?price_range=0,100"}, &Pagination{30, 1}},
		{"limit=10, page=5", args{"https://example.com/fruits?price_range=0,100&page=5&limit=10"}, &Pagination{10, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePagination(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}
