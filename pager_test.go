package pager

import (
	"reflect"
	"testing"
)

func TestPager_GetActiveAndSidesLimit(t *testing.T) {
	type fields struct {
		Limit           int
		Page            int
		SidePagingCount int
		totalCount      int
	}

	tests := []struct {
		name      string
		fields    fields
		wantLimit int
		wantPage  int
	}{
		// TODO: Add test cases.
		{"limit=5&page=1 in 100 record", fields{Limit: 5, Page: 1, SidePagingCount: 2, totalCount: 100}, 25, 0},
		{"limit=2&page=5 in 100 record", fields{Limit: 2, Page: 5, SidePagingCount: 2, totalCount: 100}, 10, 2 * 2},
		{"limit=5&page=6 in 100 record", fields{Limit: 5, Page: 6, SidePagingCount: 2, totalCount: 58}, 25, 3 * 5},
		{"limit=1&page=1 in 1 record", fields{Limit: 1, Page: 1, SidePagingCount: 2, totalCount: 1}, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pager{
				limit:           tt.fields.Limit,
				page:            tt.fields.Page,
				sidePagingCount: tt.fields.SidePagingCount,
				totalCount:      tt.fields.totalCount,
			}
			gotLimit, gotPage := p.GetActiveAndSidesLimit()
			if gotLimit != tt.wantLimit {
				t.Errorf("Pager.GetActiveAndSidesLimit() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if gotPage != tt.wantPage {
				t.Errorf("Pager.GetActiveAndSidesLimit() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
		})
	}
}

func TestPager_LastPage(t *testing.T) {
	type fields struct {
		Limit      int
		totalCount int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{"limit=10 in 100 record", fields{Limit: 10, totalCount: 100}, 9},
		{"limit=2 in 100 record", fields{Limit: 2, totalCount: 100}, 49},
		{"limit=5 in 59 record", fields{Limit: 5, totalCount: 59}, 11},
		{"limit=5 in 5 record", fields{Limit: 5, totalCount: 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pager{
				limit:      tt.fields.Limit,
				totalCount: tt.fields.totalCount,
			}
			if got := p.LastPageIndex(); got != tt.want {
				t.Errorf("Pager.LastPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_StartPage(t *testing.T) {
	type fields struct {
		Limit           int
		Page            int
		SidePagingCount int
		totalCount      int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{"limit=10", fields{Limit: 10, Page: 1, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 2, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 3, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 4, SidePagingCount: 2, totalCount: 100}, 1},
		{"limit=5", fields{Limit: 5, Page: 6, SidePagingCount: 2, totalCount: 58}, 3},
		{"limit=5", fields{Limit: 5, Page: 1, SidePagingCount: 2, totalCount: 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pager{
				limit:           tt.fields.Limit,
				page:            tt.fields.Page,
				sidePagingCount: tt.fields.SidePagingCount,
				totalCount:      tt.fields.totalCount,
			}
			if got := p.StartPageIndex(); got != tt.want {
				t.Errorf("Pager.StartPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPager_GetPageCount(t *testing.T) {
	type fields struct {
		limit      int
		totalCount int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{"limit=0", fields{0, 100}, 0},
		{"limit=10, total=100 -> 10 pages", fields{10, 100}, 10},
		{"limit=10, total=101 -> 11 pages", fields{10, 101}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pager{
				limit:      tt.fields.limit,
				totalCount: tt.fields.totalCount,
			}
			if got := p.GetPageCount(); got != tt.want {
				t.Errorf("Pager.GetPageCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fruit struct {
	Name  string
	Price int
}

type fruitCondition struct {
	PriceLowerLimit  *int
	PriceHigherLimit *int
}

func newFruitCondition(low, high int) *fruitCondition {
	return &fruitCondition{
		PriceLowerLimit:  &low,
		PriceHigherLimit: &high,
	}
}

func (fc *fruitCondition) ApplyCondition(s interface{}) {
	fetcher, ok := s.(*fruitFetcher)
	if !ok {
		return
	}
	if fc.PriceHigherLimit != nil {
		fetcher.priceHigherLimit = *fc.PriceHigherLimit
	}
	if fc.PriceLowerLimit != nil {
		fetcher.priceLowerLimit = *fc.PriceLowerLimit
	}
}

type fruitFetcher struct {
	priceLowerLimit  int
	priceHigherLimit int
}

func newFruitFetcher() *fruitFetcher {
	return &fruitFetcher{
		priceLowerLimit:  -1 << 31,
		priceHigherLimit: 1<<31 - 1,
	}
}

func (ff *fruitFetcher) Count(cond ConditionApplier) (int, error) {
	if cond != nil {
		cond.ApplyCondition(ff)
	}
	dummyFruits := ff.GetDummy()
	return len(dummyFruits), nil
}

func (ff *fruitFetcher) FetchPage(limit, offset int, cond ConditionApplier, order Order, result *PageFetchResult) error {
	if cond != nil {
		cond.ApplyCondition(ff)
	}
	dummyFruits := ff.GetDummy()
	var toIndex int
	toIndex = offset + limit
	if toIndex > len(dummyFruits) {
		toIndex = len(dummyFruits)
	}
	for _, fruit := range dummyFruits[offset:toIndex] {
		*result = append(*result, fruit)
	}
	return nil
}

func (ff *fruitFetcher) GetDummy() []fruit {
	result := make([]fruit, 0)
	for _, f := range dummyFruits {
		if ff.priceHigherLimit >= f.Price && f.Price >= ff.priceLowerLimit {
			result = append(result, f)
		}
	}
	return result
}

var dummyFruits = []fruit{
	fruit{"Apple", 112},
	fruit{"Pear", 245},
	fruit{"Banana", 60},
	fruit{"Orange", 80},
	fruit{"Kiwi", 106},
	fruit{"Strawberry", 350},
	fruit{"Grape", 400},
	fruit{"Grapefruit", 150},
	fruit{"Pineapple", 200},
	fruit{"Cherry", 140},
	fruit{"Mango", 199},
}

func refInt(n int) *int {
	return &n
}
func TestGetPaging(t *testing.T) {
	type args struct {
		fetcher PagingFetcher
		setting *Setting
	}
	tests := []struct {
		name           string
		args           args
		wantTotalCount int
		wantPageCount  int
		wantRes        *PagingResponse
		wantErr        bool
	}{
		// TODO: Add test cases.
		{"active is out of range", args{newFruitFetcher(), &Setting{
			Limit:      refInt(2),
			ActivePage: refInt(100),
		}}, 0, 0, nil, true},
		{"no response", args{newFruitFetcher(), &Setting{
			Limit:      refInt(2),
			ActivePage: refInt(1),
			Cond:       newFruitCondition(-1, -1),
		}}, 0, 0, &PagingResponse{
			Pages: Pages{
				"active":         PageFetchResult{},
				"first":          PageFetchResult{},
				"last":           PageFetchResult{},
				"before_distant": nil,
				"before_near":    nil,
				"after_near":     nil,
				"after_distant":  nil,
			},
		}, false},
		{"no condition", args{newFruitFetcher(), &Setting{
			Limit:      refInt(2),
			ActivePage: refInt(1),
		}}, 11, 6,
			&PagingResponse{
				Pages: Pages{
					"active": PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
					},
					"first": PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
					},
					"last": PageFetchResult{
						dummyFruits[10],
					},
					"before_distant": PageFetchResult{
						dummyFruits[2],
						dummyFruits[3],
					},
					"before_near": PageFetchResult{
						dummyFruits[4],
						dummyFruits[5],
					},
					"after_near": PageFetchResult{
						dummyFruits[6],
						dummyFruits[7],
					},
					"after_distant": PageFetchResult{
						dummyFruits[8],
						dummyFruits[9],
					},
				},
			},
			false,
		},
		{"price 100-300", args{newFruitFetcher(), &Setting{
			Limit:      refInt(1),
			ActivePage: refInt(4),
			Cond:       newFruitCondition(100, 300),
		}}, 7, 7,
			&PagingResponse{
				Pages: Pages{
					"active": PageFetchResult{
						dummyFruits[7],
					},
					"first": PageFetchResult{
						dummyFruits[0],
					},
					"last": PageFetchResult{
						dummyFruits[10],
					},
					"before_distant": PageFetchResult{
						dummyFruits[1],
					},
					"before_near": PageFetchResult{
						dummyFruits[4],
					},
					"after_near": PageFetchResult{
						dummyFruits[8],
					},
					"after_distant": PageFetchResult{
						dummyFruits[9],
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotalCount, gotPageCount, gotRes, err := GetPaging(tt.args.fetcher, tt.args.setting)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaging() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotalCount != tt.wantTotalCount {
				t.Errorf("GetPaging() gotTotalCount = %v, want %v", gotTotalCount, tt.wantTotalCount)
			}
			if gotPageCount != tt.wantPageCount {
				t.Errorf("GetPaging() gotPageCount = %v, want %v", gotPageCount, tt.wantPageCount)
			}
			if err != nil {
				return
			}
			for key := range tt.wantRes.Pages {
				gotVal, ok := gotRes.Pages[key]
				if !ok {
					t.Errorf("GetPaging() gotRes.Pages must have key = %v, but got %v", key, gotRes)
				}
				if !reflect.DeepEqual(tt.wantRes.Pages[key], gotVal) {
					t.Errorf("GetPaging() gotRes.Pages[%v] = %v, want %v", key, gotVal, tt.wantRes.Pages[key])
				}
			}
		})
	}
}
