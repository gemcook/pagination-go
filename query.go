package pagination

import (
	"net/url"
	"strconv"
)

// Query has pagination query parameters.
type Query struct {
	Limit   int
	Page    int
	Sort    []*Order
	Enabled bool
}

// ParseQuery parses URL query string to get limit, page and sort
func ParseQuery(queryStr string) *Query {

	// Set default values.
	p := &Query{Limit: 30, Page: 1, Enabled: true}

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

	if pageStr := query.Get("pagination"); pageStr != "" {
		if pageStr == "false" {
			p.Enabled = false
		}
	}

	p.Sort = ParseSort(queryStr)
	return p
}
