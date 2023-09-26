// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"bytes"
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

		result.HttpResponse.Body = io.NopCloser(bytes.NewReader(respBody))

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

		// Automation@2022-08-08 - Runbooks - returns a 200 OK with no Body
		if result.HttpResponse.StatusCode == http.StatusOK && result.HttpResponse.ContentLength == 0 {
			result.Status = pollers.PollingStatusSucceeded
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
			// ContainerRegistry@2019-06-01-preview returns `Creating` rather than `InProgress` during creation
			"Creating": pollers.PollingStatusInProgress,
			// SignalR@2022-02-01 returns `Running` rather than `InProgress` during creation
			"Running": pollers.PollingStatusInProgress,
			// KubernetesConfiguration@2022-11-01 returns `Updating` rather than `InProgress` during update
			"Updating": pollers.PollingStatusInProgress,
			// StorageSync@2020-03-01 returns `validateInput`, `newPrivateDnsEntries`, `finishNewStorageSyncService` rather than `InProgress` during creation/update
			// See: https://github.com/hashicorp/go-azure-sdk/issues/565
			"validateInput":                    pollers.PollingStatusInProgress,
			"newPrivateDnsEntries":             pollers.PollingStatusInProgress,
			"newManagedIdentityCredentialStep": pollers.PollingStatusInProgress,
			"finishNewStorageSyncService":      pollers.PollingStatusInProgress,
			// StorageSync@2020-03-01 (CloudEndpoints) returns `newReplicaGroup` rather than `InProgress` during creation/update
			// See: https://github.com/hashicorp/go-azure-sdk/issues/565
			"newReplicaGroup": pollers.PollingStatusInProgress,
			// AnalysisServices @ 2017-08-01 (Servers) returns `Provisioning` during Creation
			"Provisioning": pollers.PollingStatusInProgress,
			// Resources @ 2020-10-01 (DeploymentScripts) returns `ProvisioningResources` during Creation
			"ProvisioningResources": pollers.PollingStatusInProgress,
			// AnalysisServices @ 2017-08-01 (Servers Resume) returns `Resuming` during Update
			"Resuming": pollers.PollingStatusInProgress,
			// AnalysisServices @ 2017-08-01 (Servers Suspend) returns `Scaling` during Update
			"Scaling": pollers.PollingStatusInProgress,
			// HealthBot @ 2022-08-08 (HealthBots CreateOrUpdate) returns `Working` during Creation
			"Working": pollers.PollingStatusInProgress,
		}
		for k, v := range statuses {
			if strings.EqualFold(string(op.Properties.ProvisioningState), string(k)) {
				result.Status = v
				break
			}
			if strings.EqualFold(string(op.Status), string(k)) {
				result.Status = v
				break
			}
		}

		if result.Status == pollers.PollingStatusFailed {
			lroError, parseError := parseErrorFromApiResponse(*result.HttpResponse.Response)
			if parseError != nil {
				return nil, parseError
			}

			err = pollers.PollingFailedError{
				HttpResponse: result.HttpResponse,
				Message:      lroError.Error(),
			}
		}

		if result.Status == pollers.PollingStatusCancelled {
			lroError, parseError := parseErrorFromApiResponse(*result.HttpResponse.Response)
			if parseError != nil {
				return nil, parseError
			}

			err = pollers.PollingCancelledError{
				HttpResponse: result.HttpResponse,
				Message:      lroError.Error(),
			}
		}

		if result.Status == "" {
			err = fmt.Errorf("`result.Status` was nil/empty - `op.Status` was %q / `op.Properties.ProvisioningState` was %q", string(op.Status), string(op.Properties.ProvisioningState))
		}
	}

	return
}

type operationResult struct {
	Name *string `json:"name"`
	// Some APIs (such as CosmosDbPostgreSQLCluster) return a DateTime value that doesn't match RFC3339
	// as such we're intentionally parsing this as a string (for info) rather than as a time.Time due to:
	// > parsing time "2023-08-11 01:58:30 +0000" as "2006-01-02T15:04:05Z07:00":
	// >  cannot parse " 01:58:30 +0000" as "T"
	StartTime *string `json:"startTime"`

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
