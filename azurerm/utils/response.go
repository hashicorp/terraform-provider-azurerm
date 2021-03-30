package utils

import (
	"net"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
)

func ResponseWasNotFound(resp autorest.Response) bool {
	return HTTPResponseWasStatusCode(resp.Response, http.StatusNotFound)
}

func Track2ResponseWasNotFound(err error) bool {
	if v, ok := err.(azcore.HTTPResponse); ok {
		return HTTPResponseWasStatusCode(v.RawResponse(), http.StatusNotFound)
	}
	return false
}

func ResponseWasBadRequest(resp autorest.Response) bool {
	return HTTPResponseWasStatusCode(resp.Response, http.StatusBadRequest)
}

func Track2ResponseWasBadRequest(err error) bool {
	if v, ok := err.(azcore.HTTPResponse); ok {
		return HTTPResponseWasStatusCode(v.RawResponse(), http.StatusBadRequest)
	}
	return false
}

func ResponseWasForbidden(resp autorest.Response) bool {
	return HTTPResponseWasStatusCode(resp.Response, http.StatusForbidden)
}

func Track2ResponseWasForbidden(err error) bool {
	if v, ok := err.(azcore.HTTPResponse); ok {
		return HTTPResponseWasStatusCode(v.RawResponse(), http.StatusForbidden)
	}
	return false
}

func ResponseWasConflict(resp autorest.Response) bool {
	return HTTPResponseWasStatusCode(resp.Response, http.StatusConflict)
}

func Track2ResponseWasConflict(err error) bool {
	if v, ok := err.(azcore.HTTPResponse); ok {
		return HTTPResponseWasStatusCode(v.RawResponse(), http.StatusConflict)
	}
	return false
}

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

func ResponseWasStatusCode(resp autorest.Response, statusCode int) bool { // nolint: unparam
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}

func HTTPResponseWasStatusCode(resp *http.Response, statusCode int) bool { // nolint: unparam
	if resp == nil {
		return false
	}
	return resp.StatusCode == statusCode
}
