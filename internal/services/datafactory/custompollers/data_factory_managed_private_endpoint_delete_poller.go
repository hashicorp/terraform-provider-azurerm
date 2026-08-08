// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dataFactoryManagedPrivateEndpointDeletePoller{}

type dataFactoryManagedPrivateEndpointDeletePoller struct {
	client *managedprivateendpoints.ManagedPrivateEndpointsClient
	id     managedprivateendpoints.ManagedPrivateEndpointId
}

var (
	pollingSuccess = pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		PollInterval: 1 * time.Minute,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewDataFactoryManagedPrivateEndpointDeletePoller(client *managedprivateendpoints.ManagedPrivateEndpointsClient, id managedprivateendpoints.ManagedPrivateEndpointId) *dataFactoryManagedPrivateEndpointDeletePoller {
	return &dataFactoryManagedPrivateEndpointDeletePoller{
		client: client,
		id:     id,
	}
}

func (p dataFactoryManagedPrivateEndpointDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, managedprivateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingSuccess, nil
		} else {
			return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
		}
	}

	return &pollingInProgress, nil
}
