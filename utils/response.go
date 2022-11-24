package utils

import (
	"net"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func ResponseWasNotFound(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

func ResponseWasBadRequest(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusBadRequest)
}

func ResponseWasForbidden(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusForbidden)
}

func ResponseWasConflict(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusConflict)
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

// ResponseWasBadRequestWithNotRegistered is a workaround for https://github.com/Azure/azure-rest-api-specs/issues/12759
func ResponseWasBadRequestWithNotRegistered(resp autorest.Response, err error) bool {
	e, ok := err.(autorest.DetailedError)
	if !ok {
		return false
	}

	r, ok := e.Original.(*azure.RequestError)
	if !ok {
		return false
	}

	if r.ServiceError == nil || len(r.ServiceError.Details) == 0 {
		return false
	}

	serviceCode, ok := r.ServiceError.Details[0]["code"]
	if !ok {
		return false
	}

	return ResponseWasStatusCode(resp, http.StatusBadRequest) && serviceCode == "SubscriptionIdNotRegisteredWithSrs"
}
