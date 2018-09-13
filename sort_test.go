package pagination

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

type XormSorterMock struct {
	Orders []string
}

func (x *XormSorterMock) Asc(cols ...string) {
	for _, s := range cols {
		x.Orders = append(x.Orders, "+"+s)
	}
}

func (x *XormSorterMock) Desc(cols ...string) {
	for _, s := range cols {
		x.Orders = append(x.Orders, "-"+s)
	}
}

func TestApplyXormOrders(t *testing.T) {

	type args struct {
		orders []*Order
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"", args{orders: []*Order{
			&Order{ColumnName: "col1", Direction: DirectionAsc},
			&Order{ColumnName: "col2", Direction: DirectionDesc},
		}}, []string{"+col1", "-col2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xormSorter := &XormSorterMock{
				Orders: []string{},
			}
			ApplyXormOrders(xormSorter, tt.args.orders)

			if !reflect.DeepEqual(xormSorter.Orders, tt.want) {
				t.Errorf("wrong result. got=%+v, want=%+v", xormSorter.Orders, tt.want)
			}
		})
	}
}
