// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type BaseClient interface {
	// Execute invokes a non-paginated API request and returns a populated *Response
	Execute(ctx context.Context, req *Request) (*Response, error)

	// ExecutePaged invokes a paginated API request, merges the results from all pages and returns a populated *Response with all results
	ExecutePaged(ctx context.Context, req *Request) (*Response, error)

	// NewRequest constructs a *Request that can be passed to Execute or ExecutePaged
	NewRequest(ctx context.Context, input RequestOptions) (*Request, error)
}

// RequestRetryFunc is a function that determines whether an HTTP request has failed due to eventual consistency and should be retried
type RequestRetryFunc func(*http.Response, *odata.OData) (bool, error)

// RequestMiddleware can manipulate or log a request before it is sent
type RequestMiddleware func(*http.Request) (*http.Request, error)

// ResponseMiddleware can manipulate or log a response before it is parsed and returned
type ResponseMiddleware func(*http.Request, *http.Response) (*http.Response, error)

// ValidStatusFunc is a function that tests whether an HTTP response is considered valid for the particular request.
type ValidStatusFunc func(*http.Response, *odata.OData) bool
