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

### parse Function

Package `pagination` provides `ParseQuery` and `ParseMap` functions that parses Query Parameters from request URL.
Those query parameters below will be parsed.

| query parameter | Mapped field | required | expected value | default value |
| --- | --- | --- | --- | --- |
| `limit` | `Limit` | no | positive integer | `30` |
| `page` | `Page` | no | positive integer (1~) | `1` |
| `pagination` | `Enabled` | no | boolean | `true` |

#### Query String from URL

```go
// RequestURI: https://example.com/fruits?limit=10&page=1&price_range=100,300&sort=+price&pagination=true
p := pagination.ParseQuery(r.URL.RequestURI())

fmt.Println("limit =", p.Limit)
fmt.Println("page =", p.Page)
fmt.Println("pagination =", p.Enabled)
```

#### Query Parameters from AWS API Gateway - Lambda

```go
import "github.com/aws/aws-lambda-go/events"

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// event.QueryStringParameters
	// map[string]string{"limit": "10", "page": "1", "pagination": "false"}

	p := pagination.ParseMap(event.QueryStringParameters)
	fmt.Println("limit =", p.Limit)
	fmt.Println("page =", p.Page)
	fmt.Println("pagination =", p.Enabled)
}
```

### fetching condition [OPTIONAL]

Tell pagination the condition to filter resources.
Then use `cond.ApplyCondition` in `Count` and `FetchPage` function.
`ApplyCondition` takes a single parameter to pass your resource dependent object (something like O/R mapper).

### Orders [OPTIONAL]

Optionally, pagination takes orders.
Use `pagination.ParseQuery` or `pagination.ParseMap` to parse sort parameter in query string.
Then, just pass `Query.Sort` to `Setting.Orders`.

Those query parameters below will be parsed.

| query parameter | Mapped field | required | expected value | default value |
| --- | --- | --- | --- | --- |
| `sort` | `Sort` | no | `+column_name` for ascending sort. </br> `-column_name` for descending sort. | `nil` |

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
