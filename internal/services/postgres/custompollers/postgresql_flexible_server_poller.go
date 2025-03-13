// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

var _ pollers.PollerType = &postgresqlFlexibleServerPoller{}

type postgresqlFlexibleServerPoller struct {
	client   *resourcemanager.Client
	response *http.Response
}

type operationResult struct {
	Name      *string               `json:"name"`
	StartTime *string               `json:"startTime"`
	Status    pollers.PollingStatus `json:"status"`
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func postgresqlFlexibleServerPollingUriForLongRunningOperation(resp *http.Response) string {
	pollingUrl := resp.Header.Get("Azure-AsyncOperation")
	if pollingUrl == "" {
		pollingUrl = resp.Header.Get("Location")
	}
	return pollingUrl
}

func NewPostgresqlFlexibleServerPoller(client *resourcemanager.Client, resp *http.Response) *postgresqlFlexibleServerPoller {
	return &postgresqlFlexibleServerPoller{
		client:   client,
		response: resp,
	}
}

func (p postgresqlFlexibleServerPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	pollingUrl := postgresqlFlexibleServerPollingUriForLongRunningOperation(p.response)
	if pollingUrl == "" {
		return nil, fmt.Errorf("no polling URL found in response")
	}

	parsedPollingUrl, err := url.Parse(pollingUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid polling URL %q in response: %v", pollingUrl, err)
	}
	if !parsedPollingUrl.IsAbs() {
		return nil, fmt.Errorf("invalid polling URL %q in response: URL was not absolute", pollingUrl)
	}

	reqOpts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusNotFound,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: nil,
		Path:          parsedPollingUrl.Path,
	}

	req, err := p.client.NewRequest(ctx, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("building request for long-running-operation: %+v", err)
	}
	req.URL.RawQuery = parsedPollingUrl.RawQuery

	resp, err := req.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("executing request: %+v", err)
	}

	if resp.Response != nil {
		var respBody []byte
		respBody, err = io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("parsing response body: %+v", err)
		}
		resp.Response.Body.Close()

		resp.Response.Body = io.NopCloser(bytes.NewReader(respBody))

		contentType := resp.Response.Header.Get("Content-Type")
		var op operationResult
		if strings.Contains(strings.ToLower(contentType), "application/json") {
			if err = json.Unmarshal(respBody, &op); err != nil {
				return nil, fmt.Errorf("unmarshalling response body: %+v", err)
			}
		} else {
			return nil, fmt.Errorf("internal-error: polling support for the Content-Type %q was not implemented: %+v", contentType, err)
		}

		if op.Status == pollers.PollingStatusInProgress {
			return &pollingInProgress, nil
		}
	}

	return &pollingSuccess, nil
}
