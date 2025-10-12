// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ pollers.PollerType = &machineLearningRegistryPoller{}

type machineLearningRegistryPoller struct {
	client     *registrymanagement.RegistryManagementClient
	id         registrymanagement.RegistryId
	pollingUrl *url.URL
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

// NewMachineLearningRegistryPoller creates a new poller for Machine Learning Registry operations
// This handles the known Azure ML Registry API bug where CreateOrUpdate returns 202 with no body instead of the 201 in the spec.
func NewMachineLearningRegistryPoller(client *registrymanagement.RegistryManagementClient, id registrymanagement.RegistryId, response *http.Response) (*machineLearningRegistryPoller, error) {
	// Only create the poller if we receive a 202 status, indicating async operation
	if response == nil {
		return nil, fmt.Errorf("no response provided")
	}
	defer response.Body.Close()
	bodyString := "[not able to read response body]"
	bodyBytes, err := io.ReadAll(response.Body)
	if err == nil {
		bodyString = string(bodyBytes)
	}
	if response.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("invalid response status: %d, body: %s", response.StatusCode, bodyString)
	}
	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("invalid response status: %d, body: %s", response.StatusCode, bodyString)
	}

	pollingUrl := response.Header.Get("Azure-AsyncOperation")
	if pollingUrl == "" {
		pollingUrl = response.Header.Get("Location")
	}

	if pollingUrl == "" {
		return nil, fmt.Errorf("no polling URL found in response (neither Azure-AsyncOperation nor Location headers were present)")
	}

	url, err := url.Parse(pollingUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid polling URL %q in response: %v", pollingUrl, err)
	}
	if !url.IsAbs() {
		return nil, fmt.Errorf("invalid polling URL %q in response: URL was not absolute", pollingUrl)
	}
	url.Query().Encode()

	return &machineLearningRegistryPoller{
		client:     client,
		pollingUrl: url,
	}, nil
}

type operationResult struct {
	Status status `json:"status"`
}

type status string

var (
	pollingDeleteSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 20 * time.Second,
	}
	pollingDeleteInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 20 * time.Second,
	}
)

type myOptions struct {
	azureAsyncOperation string
}

var _ client.Options = myOptions{}

func (p myOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (p myOptions) ToOData() *odata.Query {
	return &odata.Query{}
}

func (p myOptions) ToQuery() *client.QueryParams {
	u, err := url.Parse(p.azureAsyncOperation)
	if err != nil {
		log.Printf("[ERROR] Unable to parse Azure-AsyncOperation URL: %v", err)
		return nil
	}
	q := client.QueryParams{}
	for k, v := range u.Query() {
		if len(v) > 0 {
			q.Append(k, v[0])
		}
	}
	return &q
}

func (p machineLearningRegistryPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.pollingUrl == nil {
		return nil, fmt.Errorf("internal error: cannot poll without a pollingUrl")
	}

	reqOpts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       p.pollingUrl.Path,
		OptionsObject: myOptions{
			azureAsyncOperation: p.pollingUrl.String(),
		},
	}

	req, err := p.client.Client.NewRequest(ctx, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := p.client.Client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.pollingUrl.String(), err)
	}
	var respBody struct {
		Status          string  `json:"status"` // "InProgress",  "Succeeded"
		PercentComplete float32 `json:"percentComplete"`
		registrymanagement.RegistryTrackedResource
	}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("decoding response body: %+v", err)
	}
	if respBody.Status == "InProgress" {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: 10 * time.Second,
		}, nil
	}
	if respBody.Status == "Succeeded" {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 10 * time.Second,
		}, nil
	}
	// The status URL can also respond with the full registry object
	if pointer.From(respBody.Type) == "Microsoft.MachineLearningServices/registries" {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 10 * time.Second,
		}, nil
	}

	return nil, fmt.Errorf("unexpected status code %d. Response body: %s", resp.StatusCode, resp.Body)
}
