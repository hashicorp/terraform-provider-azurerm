// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
)

type sreAgentTransport func(*http.Request) (*http.Response, error)

func (f sreAgentTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func sreAgentHTTPResponse(request *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header: http.Header{"Content-Type": []string{"application/json"}, sdkclient.SkipPollingDelayHeader: []string{"true"}},
		Body:   io.NopCloser(strings.NewReader(body)), Request: request, ContentLength: int64(len(body)),
	}
}
