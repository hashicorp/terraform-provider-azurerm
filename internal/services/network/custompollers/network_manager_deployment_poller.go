// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

var _ pollers.PollerType = &networkManagerDeploymentPoller{}

type networkManagerDeploymentPoller struct {
	client *networkmanagers.NetworkManagersClient
	id     parse.ManagerDeploymentId
}

func NewNetworkManagerDeploymentPoller(client *networkmanagers.NetworkManagersClient, id parse.ManagerDeploymentId) pollers.PollerType {
	return &networkManagerDeploymentPoller{
		client: client,
		id:     id,
	}
}

func (p networkManagerDeploymentPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result := &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}

	payload := networkmanagers.NetworkManagerDeploymentStatusParameter{
		Regions:         &[]string{location.Normalize(p.id.Location)},
		DeploymentTypes: &[]networkmanagers.ConfigurationType{networkmanagers.ConfigurationType(p.id.ScopeAccess)},
	}

	networkManagerId := networkmanagers.NewNetworkManagerID(p.id.SubscriptionId, p.id.ResourceGroup, p.id.NetworkManagerName)
	resp, err := p.client.NetworkManagerDeploymentStatusList(ctx, networkManagerId, payload)
	if err != nil {
		result.Status = pollers.PollingStatusFailed
		return result, pollers.PollingFailedError{
			Message: fmt.Errorf("listing deployments for %s", networkManagerId).Error(),
		}
	}

	// List request may initially return no results
	if resp.Model == nil || resp.Model.Value == nil || len(*resp.Model.Value) == 0 {
		return result, nil
	}

	if *(*resp.Model.Value)[0].DeploymentStatus == networkmanagers.DeploymentStatusDeployed {
		result.Status = pollers.PollingStatusSucceeded
		return result, nil
	}

	return result, nil
}
