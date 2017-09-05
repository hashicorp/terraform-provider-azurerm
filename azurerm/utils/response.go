package utils

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

func ResponseWasConflict(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusConflict)
}

func ResponseWasNotFound(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusNotFound)
}

func responseWasStatusCode(resp autorest.Response, statusCode int) bool {
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}
