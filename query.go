package rest

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Filters map[string]interface{}
	Order   []string
	Offset  int
	Limit   int
	//Embed []string
}

func (q *Query) Render() string {
	vals := url.Values{}
	for k, v := range q.Filters {
		vals.Set(k, v)
	}

	if q.Order != nil && len(q.Order) > 0 {
		vals.Set("_order", strings.Join(q.Order, ","))
	}

	if q.Offset > 0 {
		vals.Set("_offset", strconv.Itoa(q.Offset))
	}

	if q.Limit > 0 {
		vals.Set("_limit", strconv.Itoa(q.Limit))
	}

	return vals.Encode()
}
