package pagination

import (
	"fmt"
	"math"
	"strconv"
)

// Setting is pagination setting
type Setting struct {
	// data record count per single page
	Limit int `json:"limit"`
	// active page number　(1〜)
	Page   int `json:"page"`
	Cond   interface{}
	Orders []*Order
}

// Pager has pagination parameters
type Pager struct {
	limit           int
	page            int
	sidePagingCount int
	totalCount      int
	Condition       interface{}
	Orders          []*Order
	fetcher         PageFetcher
}

// PageFetcher is the interface to fetch the desired range of record.
type PageFetcher interface {
	Count(cond interface{}) (int, error)
	FetchPage(cond interface{}, input *PageFetchInput, result *PageFetchResult) error
}

// PageFetchInput input for page fetcher
type PageFetchInput struct {
	Limit  int
	Offset int
	Orders []*Order
}

// GetPageName returns named page
func GetPageName(i int) string {
	switch i {
	case 0:
		return "before_distant"
	case 1:
		return "before_near"
	case 2:
		return "after_near"
	case 3:
		return "after_distant"
	default:
		return strconv.Itoa(i)
	}
}

// Fetch returns paging response using arbitrary record fetcher.
func Fetch(fetcher PageFetcher, setting *Setting) (totalCount, pageCount int, res *PagingResponse, err error) {
	pager, err := newPager(fetcher, setting)
	if err != nil {
		return 0, 0, nil, err
	}
	res, err = pager.GetPages()
	if err != nil {
		return 0, 0, nil, err
	}

	return pager.totalCount, pager.GetPageCount(), res, nil
}

func newPager(fetcher PageFetcher, setting *Setting) (*Pager, error) {
	pager := Pager{}
	pager.init()
	pager.fetcher = fetcher

	if setting.Limit != 0 {
		pager.limit = setting.Limit
	}

	if setting.Page != 0 {
		if setting.Page < 1 {
			return nil, fmt.Errorf("page must be >= 1")
		}
		pager.page = setting.Page
	}

	// currently side pages count is fixed to 2
	pager.sidePagingCount = 2

	pager.Condition = setting.Cond
	pager.Orders = setting.Orders

	return &pager, nil
}

// init は Pager パラメータの初期値をセットする
func (p *Pager) init() {
	p.limit = 10
	p.page = 1
	p.sidePagingCount = 2
}

// ActivePageIndex はアクティブのページ番号を取得する
func (p *Pager) ActivePageIndex() int {
	return p.page - 1
}

// StartPageIndex は最初のページ番号を取得する
func (p *Pager) StartPageIndex() int {
	startPageIndex := (p.page - 1) - p.sidePagingCount

	// 最終ページを含む場合は取得開始位置を調整する
	endPageIndex := startPageIndex + (p.sidePagingCount * 2)
	if endPageIndex > p.LastPageIndex() {
		startPageIndex = startPageIndex - (endPageIndex - p.LastPageIndex())
	}

	if startPageIndex < 0 {
		startPageIndex = 0
	}

	return startPageIndex
}

// LastPageIndex returns the last page index.
func (p *Pager) LastPageIndex() int {
	if p.totalCount < 1 || p.limit == 0 {
		return 0
	}

	// calculate the last page index
	lastPageIndex := (p.totalCount - 1) / p.limit
	return lastPageIndex
}

// GetActiveAndSidesLimit gets records count and offset of pages chunk.
func (p *Pager) GetActiveAndSidesLimit() (limit, offset int) {
	// start record index of side pages chunk
	offset = p.StartPageIndex() * p.limit

	if offset > p.totalCount {
		offset = p.totalCount - 1
	}

	// data record limit of side pages chunk
	limit = ((p.sidePagingCount * 2) + 1) * p.limit

	if limit > p.totalCount {
		limit = p.totalCount
	}

	return limit, offset
}

// GetPages gets formated paging response.
func (p *Pager) GetPages() (*PagingResponse, error) {

	count, err := p.fetcher.Count(p.Condition)
	if err != nil {
		return nil, err
	}
	p.totalCount = count

	pageCount := p.GetPageCount()
	if pageCount == 0 {
		return p.formatResponse(PageFetchResult{}, PageFetchResult{}, PageFetchResult{}), nil
	}
	if p.page > pageCount {
		return nil, fmt.Errorf("page is out of range. page range is 1-%v", pageCount)
	}

	// active と sides に相当する範囲をまとめて取得する
	limit, offset := p.GetActiveAndSidesLimit()
	activeAndSides := make(PageFetchResult, 0, limit)
	fetchActiveInput := &PageFetchInput{
		Limit:  limit,
		Offset: offset,
		Orders: p.Orders,
	}
	err = p.fetcher.FetchPage(p.Condition, fetchActiveInput, &activeAndSides)
	if err != nil {
		return nil, err
	}

	// 最初のページが範囲外の場合は取得する
	first := make(PageFetchResult, 0, p.limit)
	if p.StartPageIndex() > 0 {
		fetchFirstInput := &PageFetchInput{
			Limit:  p.limit,
			Offset: 0,
			Orders: p.Orders,
		}
		err = p.fetcher.FetchPage(p.Condition, fetchFirstInput, &first)
		if err != nil {
			return nil, err
		}
	}

	// 最後のページが範囲外の場合は取得する
	last := make(PageFetchResult, 0, p.limit)
	if p.StartPageIndex()+(p.sidePagingCount*2) < p.LastPageIndex() {

		fetchLastInput := &PageFetchInput{
			Limit:  p.limit,
			Offset: p.LastPageIndex() * p.limit,
			Orders: p.Orders,
		}
		err = p.fetcher.FetchPage(p.Condition, fetchLastInput, &last)
		if err != nil {
			return nil, err
		}
	}

	return p.formatResponse(first, activeAndSides, last), nil
}

// GetPageCount はページの総数を返します
func (p *Pager) GetPageCount() int {
	if p.limit == 0 {
		return 0
	}
	count := math.Ceil(float64(p.totalCount) / float64(p.limit))
	return int(count)
}

// PageFetchResult has a single page chunk.
type PageFetchResult []interface{}

// Pages is a named map of pager.
type Pages map[string]PageFetchResult

// PagingResponse is a response of pager.
type PagingResponse struct {
	Pages Pages `json:"pages"`
}

func (p *Pager) formatResponse(first PageFetchResult, activeAndSides PageFetchResult, last PageFetchResult) *PagingResponse {
	active := make(PageFetchResult, 0)
	sidesLen := p.sidePagingCount * 2
	sides := make([]PageFetchResult, sidesLen, sidesLen)

	page := p.StartPageIndex() + 1
	pageIndex := 0
	for i, item := range activeAndSides {

		// fill the active page data
		if page == p.page {
			active = append(active, item)
		}
		// fill the side pages sequentially
		if page != p.page {
			sides[pageIndex] = append(sides[pageIndex], item)
		}

		// fill the first, if the chunk data has the first page.
		if page == 1 {
			first = append(first, item)
		}
		// fill the last, if the chunk data has the last page.
		if (p.LastPageIndex() + 1) == page {
			last = append(last, item)
		}

		// ページの区切り
		if (i+1)%p.limit == 0 {
			page++
			if pageIndex < sidesLen && len(sides[pageIndex]) > 0 {
				pageIndex++
			}
		}
	}

	// name pages
	responsePage := make(Pages)
	responsePage["active"] = active
	responsePage["first"] = first
	responsePage["last"] = last

	for i, sampleItems := range sides {
		pageName := GetPageName(i)
		responsePage[pageName] = sampleItems
	}

	return &PagingResponse{
		Pages: responsePage,
	}
}
