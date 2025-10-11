// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &machineLearningRegistryPoller{}

type machineLearningRegistryPoller struct {
	client *registrymanagement.RegistryManagementClient
	id     registrymanagement.RegistryId
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
// This handles the known Azure ML Registry API bug where CreateOrUpdate returns 202 with no body
func NewMachineLearningRegistryPoller(client *registrymanagement.RegistryManagementClient, id registrymanagement.RegistryId, response *http.Response) *machineLearningRegistryPoller {
	// Only create the poller if we receive a 202 status, indicating async operation
	if response.StatusCode != http.StatusAccepted {
		return nil
	}

	return &machineLearningRegistryPoller{
		client: client,
		id:     id,
	}
}

func (p machineLearningRegistryPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.RegistriesGet(ctx, p.id)
	if err != nil {
		// If the resource is not found, continue polling
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingInProgress, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	// If we successfully retrieved the resource, we're done
	if resp.Model != nil {
		return &pollingSuccess, nil
	}

	// Otherwise, continue polling
	return &pollingInProgress, nil
}
