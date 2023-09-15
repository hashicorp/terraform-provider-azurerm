// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"net/http"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type Options interface {
	// ToHeaders yields a custom Headers struct to be appended to the request
	ToHeaders() *Headers

	// ToOData yields a custom *odata.Query struct to be appended to the request
	ToOData() *odata.Query

	// ToQuery yields a custom *QueryParams struct to be appended to the request
	ToQuery() *QueryParams
}

// Headers is a representation of the HTTP headers to be sent with a Request
type Headers struct {
	vals map[string]string
}

// AppendHeader appends the http.Header values
func (h *Headers) AppendHeader(h2 http.Header) {
	if h.vals == nil {
		h.vals = map[string]string{}
	}
	for k, v := range h2 {
		if len(v) > 0 {
			h.vals[k] = v[0]
		}
	}
}

// Append sets a single header value
func (h *Headers) Append(key, value string) {
	if h.vals == nil {
		h.vals = map[string]string{}
	}
	h.vals[key] = value
}

// Merge copies the header values from h2, overwriting as necessary
func (h *Headers) Merge(h2 Headers) {
	if h.vals == nil {
		h.vals = map[string]string{}
	}
	if h2.vals == nil {
		return
	}
	for k, v := range h2.vals {
		h.vals[k] = v
	}
}

// Headers returns an http.Headers map containing header values
func (h *Headers) Headers() http.Header {
	out := make(http.Header)
	for k, v := range h.vals {
		out.Add(k, v)
	}
	return out
}

// QueryParams is a representation of the URL query parameters to be sent with a Request
type QueryParams struct {
	vals map[string]string
}

// AppendValues appends the url.Values values
func (q *QueryParams) AppendValues(q2 url.Values) {
	if q.vals == nil {
		q.vals = map[string]string{}
	}
	for k, v := range q2 {
		if len(v) > 0 {
			q.vals[k] = v[0]
		}
	}
}

// Append sets a single query parameter value
func (q *QueryParams) Append(key, value string) {
	if q.vals == nil {
		q.vals = map[string]string{}
	}
	q.vals[key] = value
}

// Merge copies the query parameter values from q2, overwriting as necessary
func (q *QueryParams) Merge(q2 Headers) {
	if q.vals == nil {
		q.vals = map[string]string{}
	}
	if q2.vals == nil {
		return
	}
	for k, v := range q2.vals {
		q.vals[k] = v
	}
}

// Values returns a url.Values map containing query parameter values
func (q *QueryParams) Values() url.Values {
	va := make(url.Values)
	for k, v := range q.vals {
		va.Set(k, v)
	}
	return va
}
