// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &automationRuntimeEnvironmentPackagePoller{}

type automationRuntimeEnvironmentPackagePoller struct {
	client *packageresource.PackageResourceClient
	id     packageresource.PackageId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewAutomationRuntimeEnvironmentPackagePoller(client *packageresource.PackageResourceClient, id packageresource.PackageId) *automationRuntimeEnvironmentPackagePoller {
	return &automationRuntimeEnvironmentPackagePoller{
		client: client,
		id:     id,
	}
}

func (p automationRuntimeEnvironmentPackagePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("polling for %s: `model` was nil", p.id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("polling for %s: `properties` was nil", p.id)
	}

	props := resp.Model.Properties

	if props.Error != nil && props.Error.Message != nil && *props.Error.Message != "" {
		return nil, pollers.PollingFailedError{
			Message: *props.Error.Message,
		}
	}

	if props.ProvisioningState == nil {
		return nil, fmt.Errorf("polling for %s: `provisioningState` was nil", p.id)
	}

	switch *props.ProvisioningState {
	case packageresource.PackageProvisioningStateSucceeded:
		return &pollingSuccess, nil

	case packageresource.PackageProvisioningStateFailed:
		return nil, pollers.PollingFailedError{
			Message: fmt.Sprintf("provisioning of %s failed", p.id),
		}

	case packageresource.PackageProvisioningStateCanceled:
		return nil, pollers.PollingFailedError{
			Message: fmt.Sprintf("provisioning of %s was canceled", p.id),
		}

	default:
		return &pollingInProgress, nil
	}
}
