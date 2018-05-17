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
		// TODO: Add test cases.
		{"default limit=30, page=1", args{"https://example.com/fruits?price_range=0,100"}, &Query{30, 1, []*Order{}}},
		{"limit=10, page=5", args{"https://example.com/fruits?price_range=0,100&page=5&limit=10"}, &Query{10, 5, []*Order{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQuery(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
