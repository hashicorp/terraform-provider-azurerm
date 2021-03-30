// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import "net/http"

// NewResponseError wraps the specified error with an error that provides access to an HTTP response.
// If an HTTP request fails, wrap the response and the associated error in this error type so that
// callers can access the underlying *http.Response as required.
// You MUST supply an inner error.
func NewResponseError(inner error, resp *http.Response) error {
	return &ResponseError{inner: inner, resp: resp}
}

// ResponseError associates an error with an HTTP response.
// Exported for type assertion purposes in azcore, use NewResponseError().
type ResponseError struct {
	inner error
	resp  *http.Response
}

// Error implements the error interface for type ResponseError.
func (e *ResponseError) Error() string {
	return e.inner.Error()
}

// Unwrap returns the inner error.
func (e *ResponseError) Unwrap() error {
	return e.inner
}

// RawResponse returns the HTTP response associated with this error.
func (e *ResponseError) RawResponse() *http.Response {
	return e.resp
}
