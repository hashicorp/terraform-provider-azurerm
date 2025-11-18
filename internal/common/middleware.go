// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

func correlationRequestIDMiddleware(id string) client.RequestMiddleware {
	return func(request *http.Request) (*http.Request, error) {
		// ensure the `X-Correlation-ID` field is set
		if request.Header.Get(HeaderCorrelationRequestID) == "" {
			request.Header.Add(HeaderCorrelationRequestID, id)
		}
		return request, nil
	}
}

func requestLoggerMiddleware(providerName string) client.RequestMiddleware {
	return func(request *http.Request) (*http.Request, error) {
		// strip the authorization header prior to printing
		authHeaderName := "Authorization"
		auth := request.Header.Get(authHeaderName)
		if auth != "" {
			request.Header.Del(authHeaderName)
		}

		// dump request to wire format
		if dump, err := httputil.DumpRequestOut(request, true); err == nil {
			log.Printf("[DEBUG] %s Request: \n%s\n", providerName, dump)
		} else {
			// fallback to basic message
			log.Printf("[DEBUG] %s Request: %s to %s\n", providerName, request.Method, request.URL)
		}

		// add the auth header back
		if auth != "" {
			request.Header.Add(authHeaderName, auth)
		}

		return request, nil
	}
}

func responseLoggerMiddleware(providerName string) client.ResponseMiddleware {
	return func(request *http.Request, response *http.Response) (*http.Response, error) {
		// dump response to wire format
		if dump, err2 := httputil.DumpResponse(response, true); err2 == nil {
			log.Printf("[DEBUG] %s Response for %s: \n%s\n", providerName, request.URL, dump)
		} else {
			// fallback to basic message
			log.Printf("[DEBUG] %s Response: %s for %s\n", providerName, response.Status, request.URL)
		}
		return response, nil
	}
}
