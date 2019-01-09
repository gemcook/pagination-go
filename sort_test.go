package pagination_test

import (
	"reflect"
	"testing"

	pagination "github.com/gemcook/pagination-go"
)

func TestParseSort(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want []*pagination.Order
	}{
		{"no sort", args{""}, []*pagination.Order{}},
		{"single col asc", args{"?sort=+col_a"}, []*pagination.Order{&pagination.Order{Direction: pagination.DirectionAsc, ColumnName: "col_a"}}},
		{"single col desc", args{"?sort=-col_b"}, []*pagination.Order{&pagination.Order{Direction: pagination.DirectionDesc, ColumnName: "col_b"}}},
		{"multi col", args{"?sort=+col_c-col_d"}, []*pagination.Order{
			&pagination.Order{Direction: pagination.DirectionAsc, ColumnName: "col_c"},
			&pagination.Order{Direction: pagination.DirectionDesc, ColumnName: "col_d"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pagination.ParseSort(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
