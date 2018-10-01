package pagination

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want *Query
	}{
		{"default", args{"https://example.com/fruits?price_range=0,100"}, &Query{30, 1, []*Order{}, true}},
		{"limit=10, page=5", args{"https://example.com/fruits?price_range=0,100&page=5&limit=10"}, &Query{10, 5, []*Order{}, true}},
		{"pagination disabled", args{"https://example.com/fruits?price_range=0,100&pagination=false"}, &Query{30, 1, []*Order{}, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQuery(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMap(t *testing.T) {
	tests := []struct {
		name string
		qs   map[string]string
		want *Query
	}{
		{"default", map[string]string{}, &Query{Limit: 30, Page: 1, Sort: []*Order{}, Enabled: true}},
		{"limit=10, page=2", map[string]string{"limit": "10", "page": "2"}, &Query{Limit: 10, Page: 2, Sort: []*Order{}, Enabled: true}},
		{"pagination=false", map[string]string{"pagination": "false"}, &Query{Limit: 30, Page: 1, Sort: []*Order{}, Enabled: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMap(tt.qs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
