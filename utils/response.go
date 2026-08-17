// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package utils is deprecated and has been mostly refactored into the helpers package.
// Deprecated: Use the helpers package instead.
package utils

import (
	"net"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

// ResponseWasNotFound is deprecated and should not be used in new code
// Deprecated
func ResponseWasNotFound(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

// ResponseWasBadRequest is deprecated and should not be used in new code
// Deprecated
func ResponseWasBadRequest(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusBadRequest)
}

// ResponseWasForbidden is deprecated and should not be used in new code
// Deprecated
func ResponseWasForbidden(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusForbidden)
}

// ResponseWasConflict is deprecated and should not be used in new code
// Deprecated
func ResponseWasConflict(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusConflict)
}

// ResponseErrorIsRetryable is deprecated and should not be used in new code
// Deprecated
func ResponseErrorIsRetryable(err error) bool {
	if arerr, ok := err.(autorest.DetailedError); ok {
		err = arerr.Original
	}

	// nolint gocritic
	switch e := err.(type) {
	case net.Error:
		if e.Temporary() || e.Timeout() {
			return true
		}
	}

	return false
}

// ResponseWasStatusCode is deprecated and should not be used in new code
// Deprecated
func ResponseWasStatusCode(resp autorest.Response, statusCode int) bool {
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}
