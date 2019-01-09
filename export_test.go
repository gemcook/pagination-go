package pagination

var CreateMockPager = func(limit, page, sidePagingCount, totalCount int) *Pager {
	pager := Pager{
		limit:           limit,
		page:            page,
		sidePagingCount: sidePagingCount,
		totalCount:      totalCount,
	}
	return &pager
}

var NewPager = newPager
