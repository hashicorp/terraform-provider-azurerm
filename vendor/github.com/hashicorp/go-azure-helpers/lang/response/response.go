// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package response

import (
	"net/http"
)

// WasBadRequest returns true if the HttpResponse is non-nil and has a status code of BadRequest
func WasBadRequest(resp *http.Response) bool {
	return WasStatusCode(resp, http.StatusBadRequest)
}

// WasConflict returns true if the HttpResponse is non-nil and has a status code of Conflict
func WasConflict(resp *http.Response) bool {
	return WasStatusCode(resp, http.StatusConflict)
}

// WasForbidden returns true if the HttpResponse is non-nil and has a status code of Forbidden
func WasForbidden(resp *http.Response) bool {
	return WasStatusCode(resp, http.StatusForbidden)
}

// WasNotFound returns true if the HttpResponse is non-nil and has a status code of NotFound
func WasNotFound(resp *http.Response) bool {
	return WasStatusCode(resp, http.StatusNotFound)
}

// WasStatusCode returns true if the HttpResponse is non-nil and matches the Status Code
// It's recommended to use WasBadRequest/WasConflict/WasNotFound where possible instead
func WasStatusCode(resp *http.Response, statusCode int) bool {
	if r := resp; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}

// WasStatusCodes returns true if the HttpResponse is non-nil and matches one of the provided Status Codes
// It's recommended to use WasBadRequest/WasConflict/WasNotFound where possible instead
func WasStatusCodes(resp *http.Response, statusCodes ...int) (result bool) {
	if resp == nil {
		return
	}

	for _, c := range statusCodes {
		if ok := WasStatusCode(resp, c); ok {
			return true
		}
	}

	return
}
