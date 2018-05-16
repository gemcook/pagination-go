package pagination

import (
	"reflect"
	"testing"
)

func TestParseQueryString(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want *QueryString
	}{
		// TODO: Add test cases.
		{"default limit=30, page=1", args{"https://example.com/fruits?price_range=0,100"}, &QueryString{30, 1}},
		{"limit=10, page=5", args{"https://example.com/fruits?price_range=0,100&page=5&limit=10"}, &QueryString{10, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQueryString(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}
