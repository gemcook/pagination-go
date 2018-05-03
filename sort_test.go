package pager

import (
	"reflect"
	"testing"
)

func TestParseSort(t *testing.T) {
	type args struct {
		queryStr string
	}
	tests := []struct {
		name string
		args args
		want []*Order
	}{
		// TODO: Add test cases.
		{"no sort", args{""}, []*Order{}},
		{"single col asc", args{"?sort=+col_a"}, []*Order{&Order{Direction: DirectionAsc, ColumnName: "col_a"}}},
		{"single col desc", args{"?sort=-col_b"}, []*Order{&Order{Direction: DirectionDesc, ColumnName: "col_b"}}},
		{"multi col", args{"?sort=+col_c-col_d"}, []*Order{
			&Order{Direction: DirectionAsc, ColumnName: "col_c"},
			&Order{Direction: DirectionDesc, ColumnName: "col_d"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSort(tt.args.queryStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
