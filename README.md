# pagination-go

[![CircleCI](https://circleci.com/gh/gemcook/pagination-go/tree/master.svg?style=shield)](https://circleci.com/gh/gemcook/pagination-go/tree/master) [![Coverage Status](https://coveralls.io/repos/github/gemcook/pagination-go/badge.svg?branch=master)](https://coveralls.io/github/gemcook/pagination-go?branch=master)

This is a helper library which perfectly matches for server-side implementation of [@gemcook/pagination](https://github.com/gemcook/pagination)

## Installation

```sh
go get -u github.com/gemcook/pagination-go
```

If you use `dep`

```sh
dep ensure -add github.com/gemcook/pagination-go
```

## Usage

### fetcher interface

First, your code to some resources must implement fetcher interface.

```go
type PageFetcher interface {
    Count(cond ConditionApplier) (int, error)
    FetchPage(limit, offset int, cond ConditionApplier, orders []*Order, result *PageFetchResult) error
}
```

### Limit and ActivePage

Use `pagination.ParseQuery` to parse query string.
Set limit and page on `Setting.Orders`.

### fetching condition [OPTIONAL]

Tell pager the condition to filter resources.
Then use `cond.ApplyCondition` in `Count` and `FetchPage` function.
`ApplyCondition` takes a single parameter to pass your resource dependent object (something like O/R mapper).

### Orders [OPTIONAL]

Optionally, pager takes orders.
Use `pagination.ParseQuery` to parse sort parameter in query string.
Then, just pass `Query.Sort` to `Setting.Orders`.

## Example

```go

import (
    "http/net"
    "encoding/json"
    "strconv"

    "github.com/gemcook/pagination-go"
)

type fruitFetcher struct{}

type FruitCondition struct{
    PriceLowerLimit int
    PriceHigherLimit int
}

func ParseFruitCondition(uri string) *FruitCondition {
    // parse uri and initialize struct
}

func (fc *FruitCondition) ApplyCondition(s interface{}) {
    // apply condition to s
}

func handler(w http.ResponseWriter, r *http.Request) {
	// RequestURI: https://example.com/fruits?limit=10&page=1&price_range=100,300&sort=+price
	p := pagination.ParseQuery(r.URL.RequestURI())
	cond := parseFruitCondition(r.URL.RequestURI())
	fetcher := newFruitFetcher()

	totalCount, totalPages, res, err := pagination.Fetch(fetcher, &pagination.Setting{
		Limit:      &p.Limit,
		ActivePage: &p.Page,
		Cond:       cond,
		Orders:     p.Sort,
	})

	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(400)
		fmt.Fprintf(w, "something wrong: %v", err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(totalCount))
	w.Header().Set("X-Total-Pages", strconv.Itoa(totalPages))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count,X-Total-Pages")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	resJSON, _ := json.Marshal(res)
	w.Write(resJSON)
}
```

For full source code, see [example/server.go](./example/server.go).

Run example.

```sh
cd example
go run server.go
```

Then open `http://localhost:8080/fruits?limit=2&page=1&price_range=100,300&sort=+price`
