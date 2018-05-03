package pager

import (
	"net/url"
	"strings"
)

// Direction shows the sort direction
type Direction string

const (
	// DirectionAsc sorts by ascending order
	DirectionAsc Direction = "ASC"
	// DirectionDesc sorts by descending order
	DirectionDesc Direction = "DESC"
)

// Order defines sort order clause
type Order struct {
	Direction  Direction
	ColumnName string
}

// ParseSort parses sort option in the given URL query string
func ParseSort(queryStr string) []*Order {
	u, err := url.Parse(queryStr)
	if err != nil {
		return []*Order{}
	}
	query := u.Query()

	if s := query.Get("sort"); s != "" {
		return ParseOrders(s)
	}
	return []*Order{}
}

// ParseOrders parses sort option string
// Sort option would be like '-col_first+col_second'.
func ParseOrders(sort string) []*Order {
	if sort == "" {
		return []*Order{}
	}

	orders := make([]*Order, 0)
	o := sort
	for _i := strings.IndexAny(o, "+- "); _i == 0; {
		col := ""
		_o := ""
		// 次に + or - が現れる位置を判定
		if i := strings.IndexAny(o[1:], "+- "); i == -1 {
			col = o[1:]
		} else {
			col = o[1 : i+1]
			_o = o[i+1:]
		}

		// カラム名が空の場合はループを抜ける
		if col == "" {
			break
		}

		// ソート条件を設定する
		var d Direction
		if o[0] == '+' || o[0] == ' ' {
			d = DirectionAsc
		} else if o[0] == '-' {
			d = DirectionDesc
		}

		orders = append(orders, &Order{d, col})

		if _o == "" {
			break
		}
		o = _o
	}
	return orders
}
