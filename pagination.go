package pager

import (
	"net/url"
	"strconv"
)

// Pagination has limit and page
type Pagination struct {
	Limit int
	Page  int
}

// ParsePagination parses URL QueryString to get limit and page
func ParsePagination(queryStr string) *Pagination {
	p := &Pagination{Limit: 30, Page: 1}

	u, err := url.Parse(queryStr)
	if err != nil {
		return p
	}
	query := u.Query()

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			p.Limit = limit
		}
	}
	if pageStr := query.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			p.Page = page
		}
	}
	return p
}
