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
