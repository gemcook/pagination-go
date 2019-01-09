package pagination_test

import (
	"reflect"
	"testing"

	pagination "github.com/gemcook/pagination-go"
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
		{"limit=5&page=1 in 100 record", fields{Limit: 5, Page: 1, SidePagingCount: 2, totalCount: 100}, 25, 0},
		{"limit=2&page=5 in 100 record", fields{Limit: 2, Page: 5, SidePagingCount: 2, totalCount: 100}, 10, 2 * 2},
		{"limit=5&page=6 in 100 record", fields{Limit: 5, Page: 6, SidePagingCount: 2, totalCount: 58}, 25, 3 * 5},
		{"limit=1&page=1 in 1 record", fields{Limit: 1, Page: 1, SidePagingCount: 2, totalCount: 1}, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pagination.CreateMockPager(
				tt.fields.Limit,
				tt.fields.Page,
				tt.fields.SidePagingCount,
				tt.fields.totalCount,
			)
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
		{"limit=10 in 100 record", fields{Limit: 10, totalCount: 100}, 9},
		{"limit=2 in 100 record", fields{Limit: 2, totalCount: 100}, 49},
		{"limit=5 in 59 record", fields{Limit: 5, totalCount: 59}, 11},
		{"limit=5 in 5 record", fields{Limit: 5, totalCount: 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pagination.CreateMockPager(
				tt.fields.Limit,
				1,
				2,
				tt.fields.totalCount,
			)
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
		{"limit=10", fields{Limit: 10, Page: 1, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 2, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 3, SidePagingCount: 2, totalCount: 100}, 0},
		{"limit=10", fields{Limit: 10, Page: 4, SidePagingCount: 2, totalCount: 100}, 1},
		{"limit=5", fields{Limit: 5, Page: 6, SidePagingCount: 2, totalCount: 58}, 3},
		{"limit=5", fields{Limit: 5, Page: 1, SidePagingCount: 2, totalCount: 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pagination.CreateMockPager(
				tt.fields.Limit,
				tt.fields.Page,
				tt.fields.SidePagingCount,
				tt.fields.totalCount,
			)
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
		{"limit=0", fields{0, 100}, 0},
		{"limit=10, total=100 -> 10 pages", fields{10, 100}, 10},
		{"limit=10, total=101 -> 11 pages", fields{10, 101}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pagination.CreateMockPager(
				tt.fields.limit,
				1,
				2,
				tt.fields.totalCount,
			)
			if got := p.GetPageCount(); got != tt.want {
				t.Errorf("Pager.GetPageCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func refInt(n int) *int {
	return &n
}
func TestFetch(t *testing.T) {
	type args struct {
		fetcher pagination.PageFetcher
		setting *pagination.Setting
	}
	tests := []struct {
		name           string
		args           args
		wantTotalCount int
		wantPageCount  int
		wantRes        *pagination.PagingResponse
		wantErr        bool
	}{
		{"active is out of range", args{newFruitFetcher(), &pagination.Setting{
			Limit:      refInt(2),
			ActivePage: refInt(100),
		}}, 0, 0, nil, true},
		{"no response", args{newFruitFetcher(), &pagination.Setting{
			Limit:      refInt(2),
			ActivePage: refInt(1),
			Cond:       newFruitCondition(-1, -1),
		}}, 0, 0, &pagination.PagingResponse{
			Pages: pagination.Pages{
				"active":         pagination.PageFetchResult{},
				"first":          pagination.PageFetchResult{},
				"last":           pagination.PageFetchResult{},
				"before_distant": nil,
				"before_near":    nil,
				"after_near":     nil,
				"after_distant":  nil,
			},
		}, false},
		{"no condition", args{newFruitFetcher(), &pagination.Setting{
			Limit:      refInt(2),
			ActivePage: refInt(1),
		}}, 11, 6,
			&pagination.PagingResponse{
				Pages: pagination.Pages{
					"active": pagination.PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
					},
					"first": pagination.PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
					},
					"last": pagination.PageFetchResult{
						dummyFruits[10],
					},
					"before_distant": pagination.PageFetchResult{
						dummyFruits[2],
						dummyFruits[3],
					},
					"before_near": pagination.PageFetchResult{
						dummyFruits[4],
						dummyFruits[5],
					},
					"after_near": pagination.PageFetchResult{
						dummyFruits[6],
						dummyFruits[7],
					},
					"after_distant": pagination.PageFetchResult{
						dummyFruits[8],
						dummyFruits[9],
					},
				},
			},
			false,
		},
		{"price 100-300", args{newFruitFetcher(), &pagination.Setting{
			Limit:      refInt(1),
			ActivePage: refInt(4),
			Cond:       newFruitCondition(100, 300),
		}}, 7, 7,
			&pagination.PagingResponse{
				Pages: pagination.Pages{
					"active": pagination.PageFetchResult{
						dummyFruits[7],
					},
					"first": pagination.PageFetchResult{
						dummyFruits[0],
					},
					"last": pagination.PageFetchResult{
						dummyFruits[10],
					},
					"before_distant": pagination.PageFetchResult{
						dummyFruits[1],
					},
					"before_near": pagination.PageFetchResult{
						dummyFruits[4],
					},
					"after_near": pagination.PageFetchResult{
						dummyFruits[8],
					},
					"after_distant": pagination.PageFetchResult{
						dummyFruits[9],
					},
				},
			},
			false,
		},
		{"totalCount is even number", args{newFruitFetcher(), &pagination.Setting{
			Limit:      refInt(5),
			ActivePage: refInt(1),
			Cond:       newFruitCondition(0, 360),
		}}, 10, 2,
			&pagination.PagingResponse{
				Pages: pagination.Pages{
					"active": pagination.PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
						dummyFruits[2],
						dummyFruits[3],
						dummyFruits[4],
					},
					"first": pagination.PageFetchResult{
						dummyFruits[0],
						dummyFruits[1],
						dummyFruits[2],
						dummyFruits[3],
						dummyFruits[4],
					},
					"last": pagination.PageFetchResult{
						dummyFruits[5],
						dummyFruits[7],
						dummyFruits[8],
						dummyFruits[9],
						dummyFruits[10],
					},
					"before_distant": pagination.PageFetchResult{
						dummyFruits[5],
						dummyFruits[7],
						dummyFruits[8],
						dummyFruits[9],
						dummyFruits[10],
					},
					"before_near":   nil,
					"after_near":    nil,
					"after_distant": nil,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotalCount, gotPageCount, gotRes, err := pagination.Fetch(tt.args.fetcher, tt.args.setting)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotalCount != tt.wantTotalCount {
				t.Errorf("Fetch() gotTotalCount = %v, want %v", gotTotalCount, tt.wantTotalCount)
			}
			if gotPageCount != tt.wantPageCount {
				t.Errorf("Fetch() gotPageCount = %v, want %v", gotPageCount, tt.wantPageCount)
			}
			if err != nil {
				return
			}
			for key := range tt.wantRes.Pages {
				gotVal, ok := gotRes.Pages[key]
				if !ok {
					t.Errorf("Fetch() gotRes.Pages must have key = %v, but got %v", key, gotRes)
				}
				if !reflect.DeepEqual(tt.wantRes.Pages[key], gotVal) {
					t.Errorf("Fetch() gotRes.Pages[%v] = %v, want %v", key, gotVal, tt.wantRes.Pages[key])
				}
			}
		})
	}
}

func TestPager_GetPages_LargeData_Last(t *testing.T) {
	type args struct {
		Limit      int
		ActivePage int
		Condition  pagination.ConditionApplier
		Orders     []*pagination.Order
		Fetcher    pagination.PageFetcher
	}
	tests := []struct {
		name     string
		args     args
		wantLast []LargeData
		wantErr  bool
	}{
		{"largeData",
			args{
				Limit: 10, ActivePage: 1, Condition: nil, Orders: []*pagination.Order{},
				Fetcher: newLargeDataFetcher(),
			},
			[]LargeData{LargeData{101}, LargeData{102}, LargeData{103}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := newLargeDataFetcher()
			p, err := pagination.NewPager(fetcher, &pagination.Setting{
				Limit:      &tt.args.Limit,
				ActivePage: &tt.args.ActivePage,
				Cond:       tt.args.Condition,
				Orders:     tt.args.Orders,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Pager.GetPages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := p.GetPages()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pager.GetPages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantLast := make(pagination.PageFetchResult, 0)
			for _, data := range tt.wantLast {
				wantLast = append(wantLast, data)
			}
			gotLast := got.Pages["last"]
			if !reflect.DeepEqual(gotLast, wantLast) {
				t.Errorf("Pager.GetPages() pages.last = %+v, wantLast %+v", gotLast, wantLast)
			}
		})
	}
}
