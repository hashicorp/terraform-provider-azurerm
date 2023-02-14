package client

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type BaseClient interface {
	Execute(ctx context.Context, req *Request) (*Response, error)
	ExecutePaged(ctx context.Context, req *Request) (*Response, error)
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
