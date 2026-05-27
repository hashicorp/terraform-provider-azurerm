// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &ManagedInstanceCreatePoller{}

const managedInstanceCreateSuccessCount = 3

type ManagedInstancesClient interface {
	Get(ctx context.Context, id commonids.SqlManagedInstanceId, options managedinstances.GetOperationOptions) (managedinstances.GetOperationResponse, error)
}

type ManagedInstanceCreatePoller struct {
	client       ManagedInstancesClient
	id           commonids.SqlManagedInstanceId
	successCount int
}

func NewManagedInstanceCreatePoller(client ManagedInstancesClient, id commonids.SqlManagedInstanceId) *ManagedInstanceCreatePoller {
	return &ManagedInstanceCreatePoller{
		client: client,
		id:     id,
	}
}

func (p *ManagedInstanceCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, managedinstances.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			p.successCount = 0
			return &pollers.PollResult{
				PollInterval: 10 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", p.id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", p.id)
	}

	switch provisioningState := pointer.From(resp.Model.Properties.ProvisioningState); provisioningState {
	case managedinstances.ProvisioningStateSucceeded:
		p.successCount++
		if p.successCount < managedInstanceCreateSuccessCount {
			return &pollers.PollResult{
				PollInterval: 10 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return &pollers.PollResult{
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil

	case managedinstances.ProvisioningStateFailed, managedinstances.ProvisioningStateCanceled:
		return nil, fmt.Errorf("creating %s: `provisioningState` was %q", p.id, provisioningState)
	}

	p.successCount = 0
	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
