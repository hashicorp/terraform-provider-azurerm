package client

import (
	"net/http"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type Options interface {
	ToHeaders() *Headers
	ToOData() *odata.Query
	ToQuery() *QueryParams
}

type Headers struct {
	val map[string]string
}

func (h Headers) Append(key, value string) {
	if h.val == nil {
		h.val = map[string]string{}
	}
	h.val[key] = value
}

func (h Headers) Headers() http.Header {
	out := make(http.Header)
	for k, v := range h.val {
		out.Add(k, v)
	}
	return out
}

type QueryParams struct {
	vals map[string]string
}

func QueryParamsFromValues(input map[string]string) *QueryParams {
	return &QueryParams{
		vals: input,
	}
}

func (q *QueryParams) Append(key, value string) {
	if q.vals == nil {
		q.vals = map[string]string{}
	}
	q.vals[key] = value
}

func (q *QueryParams) Values() url.Values {
	va := make(url.Values)
	for k, v := range q.vals {
		va.Set(k, v)
	}
	return va
}
