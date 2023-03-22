// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &longRunningOperationPoller{}

type longRunningOperationPoller struct {
	client               *client.Client
	count                int
	initialRetryDuration time.Duration
	originalUrl          *url.URL
	pollingUrl           *url.URL
}

func pollingUriForLongRunningOperation(resp *client.Response) string {
	pollingUrl := resp.Header.Get(http.CanonicalHeaderKey("Azure-AsyncOperation"))
	if pollingUrl == "" {
		pollingUrl = resp.Header.Get("Location")
	}
	return pollingUrl
}

func longRunningOperationPollerFromResponse(resp *client.Response, client *client.Client) (*longRunningOperationPoller, error) {
	poller := longRunningOperationPoller{
		client:               client,
		initialRetryDuration: 10 * time.Second,
	}

	pollingUrl := pollingUriForLongRunningOperation(resp)
	if pollingUrl == "" {
		return nil, fmt.Errorf("no polling URL found in response")
	}

	u, err := url.Parse(pollingUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid polling URL %q in response: %v", pollingUrl, err)
	}
	if !u.IsAbs() {
		return nil, fmt.Errorf("invalid polling URL %q in response: URL was not absolute", pollingUrl)
	}
	poller.pollingUrl = u

	if resp.Request != nil {
		poller.originalUrl = resp.Request.URL
	}

	if s, ok := resp.Header["Retry-After"]; ok {
		if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			poller.initialRetryDuration = time.Second * time.Duration(sleep)
		}
	}

	return &poller, nil
}

func (p *longRunningOperationPoller) Poll(ctx context.Context) (result *pollers.PollResult, err error) {
	p.count++

	if p.pollingUrl == nil {
		return nil, fmt.Errorf("internal error: cannot poll without a pollingUrl")
	}

	reqOpts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: nil,
		Path:          p.pollingUrl.Path,
	}

	// TODO: port over the `api-version` header

	req, err := p.client.NewRequest(ctx, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("building request for long-running-operation: %+v", err)
	}
	req.URL.RawQuery = p.pollingUrl.RawQuery

	// Custom RetryFunc to inspect the operation payload and check the status
	req.RetryFunc = client.RequestRetryAny(defaultRetryFunctions...)

	result = &pollers.PollResult{
		PollInterval: p.initialRetryDuration,
	}
	result.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		return nil, err
	}

	if result.HttpResponse != nil {
		var respBody []byte
		respBody, err = io.ReadAll(result.HttpResponse.Body)
		if err != nil {
			err = fmt.Errorf("parsing response body: %+v", err)
			return
		}
		result.HttpResponse.Body.Close()

		// update the poll interval if a Retry-After header is returned
		if s, ok := result.HttpResponse.Header["Retry-After"]; ok {
			if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
				result.PollInterval = time.Second * time.Duration(sleep)
			}
		}

		// 202's don't necessarily return a body, so there's nothing to deserialize
		if result.HttpResponse.StatusCode == http.StatusAccepted {
			result.Status = pollers.PollingStatusInProgress
			return
		}

		contentType := result.HttpResponse.Header.Get("Content-Type")

		var op operationResult
		if strings.Contains(strings.ToLower(contentType), "application/json") {
			if err = json.Unmarshal(respBody, &op); err != nil {
				err = fmt.Errorf("unmarshalling response body: %+v", err)
				return
			}
		} else {
			return nil, fmt.Errorf("internal-error: polling support for the Content-Type %q was not implemented: %+v", contentType, err)
		}

		if op.Properties.ProvisioningState == "" && op.Status == "" {
			return nil, fmt.Errorf("expected either `provisioningState` or `status` to be returned from the LRO API but both were empty")
		}

		// TODO: raising an error if this is Cancelled or Failed

		statuses := map[status]pollers.PollingStatus{
			statusCanceled:   pollers.PollingStatusCancelled,
			statusCancelled:  pollers.PollingStatusCancelled,
			statusFailed:     pollers.PollingStatusFailed,
			statusInProgress: pollers.PollingStatusInProgress,
			statusSucceeded:  pollers.PollingStatusSucceeded,

			// whilst the standard set above should be sufficient, some APIs differ from the spec and should be documented below:
			// Dashboard@2022-08-01 returns `Accepted` rather than `InProgress` during creation
			"Accepted": pollers.PollingStatusInProgress,
			// CostManagement@2021-10-01 returns `Completed` rather than `Succeeded`: https://github.com/Azure/azure-sdk-for-go/issues/20342
			"Completed": pollers.PollingStatusSucceeded,
			// SignalR@2022-02-01 returns `Running` rather than `InProgress` during creation
			"Running": pollers.PollingStatusInProgress,
		}
		for k, v := range statuses {
			if strings.EqualFold(string(op.Properties.ProvisioningState), string(k)) {
				result.Status = v
				return
			}
			if strings.EqualFold(string(op.Status), string(k)) {
				result.Status = v
				return
			}
		}

		if result.Status == "" {
			err = fmt.Errorf("`result.Status` was nil/empty - `op.Status` was %q / `op.Properties.ProvisioningState` was %q", string(op.Status), string(op.Properties.ProvisioningState))
		}
	}

	return
}

type operationResult struct {
	Name      *string    `json:"name"`
	StartTime *time.Time `json:"startTime"`

	Properties struct {
		// Some APIs (such as Storage) return the Resource Representation from the LRO API, as such we need to check provisioningState
		ProvisioningState status `json:"provisioningState"`
	} `json:"properties"`

	// others return Status, so we check that too
	Status status `json:"status"`
}

type status string

const (
	statusCanceled   status = "Canceled"
	statusCancelled  status = "Cancelled"
	statusFailed     status = "Failed"
	statusInProgress status = "InProgress"
	statusSucceeded  status = "Succeeded"
)
