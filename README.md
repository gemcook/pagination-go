# pager-go

[![CircleCI](https://circleci.com/gh/gemcook/pager-go/tree/master.svg?style=shield)](https://circleci.com/gh/gemcook/pager-go/tree/master)

This is a helper library which perfectly matches for server-side implementation of [@gemcook/table](https://github.com/gemcook/table)

## Installation

```sh
go get -u github.com/gemcook/pager-go
```

If you use `dep`

```sh
dep ensure -add github.com/gemcook/pager-go
```

## Usage

For the actual code. see [pager_test.go](./pager_test.go).

### fetcher interface

First, your code to some resources must implement fetcher interface.

```go
type PagingFetcher interface {
    Count(cond ConditionApplier) (int, error)
    FetchPage(limit, offset int, cond ConditionApplier, orders []*Order, result *PageFetchResult) error
}
```

### fetching condition [OPTIONAL]

Tell pager the condition to filter resources.
Then use `cond.ApplyCondition` in `Count` and `FetchPage` function.
`ApplyCondition` takes a single parameter to pass your resource dependent object (something like O/R mapper).

### Orders [OPTIONAL]

Optionally, pager takes orders.
Use `ParseSort` to parse query string sort parameter.

## Example

```go

import (
    "http/net"
    "encoding/json"
    "strconv"

    "github.com/gemcook/pager-go"
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

func Handler(w http.ResponseWriter, r *http.Request) {
    // RequestURI: https://example.com/fruits?limit=10&page=1&price_range=100,300&sort=+price
    pagination := pager.ParsePagination(r.URL.RequestURI)
    orders := pager.ParseSort(r.URL.RequestURI)
    cond := ParseFruitCondition(r.URL.RequestURI)
    fetcher := new(fruitFetcher)

    totalCount, pagesCount, res, err := pager.GetPaging(fetcher, &pager.Setting{
        Limit:      &pagination.Limit,
        ActivePage: &pagination.Page,
        Cond:       cond,
        Orders:     orders,
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
    resJson, _ := json.Marshal(res)
    w.Write(resJson)
}
```