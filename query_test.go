package pagination_test

import (
	"reflect"
	"testing"

	pagination "github.com/gemcook/pagination-go"
)

func TestParseQuery(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want *pagination.Query
	}{
		{"default", args{"https://example.com/fruits?price_range=0,100"}, &pagination.Query{Limit: 30, Page: 1, Sort: []*pagination.Order{}, Enabled: true}},
		{"limit=10, page=5", args{"https://example.com/fruits?price_range=0,100&page=5&limit=10"}, &pagination.Query{Limit: 10, Page: 5, Sort: []*pagination.Order{}, Enabled: true}},
		{"pagination disabled", args{"https://example.com/fruits?price_range=0,100&pagination=false"}, &pagination.Query{Limit: 30, Page: 1, Sort: []*pagination.Order{}, Enabled: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pagination.ParseQuery(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMap(t *testing.T) {
	tests := []struct {
		name string
		qs   map[string]string
		want *pagination.Query
	}{
		{"default", map[string]string{}, &pagination.Query{Limit: 30, Page: 1, Sort: []*pagination.Order{}, Enabled: true}},
		{"limit=10, page=2", map[string]string{"limit": "10", "page": "2"}, &pagination.Query{Limit: 10, Page: 2, Sort: []*pagination.Order{}, Enabled: true}},
		{"pagination=false", map[string]string{"pagination": "false"}, &pagination.Query{Limit: 30, Page: 1, Sort: []*pagination.Order{}, Enabled: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pagination.ParseMap(tt.qs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
