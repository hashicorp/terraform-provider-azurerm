// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"time"

	cdnFrontDoorSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ pollers.PollerType = &frontDoorOriginGroupDeletePoller{}

type frontDoorOriginGroupDeletePoller struct {
	client *cdnFrontDoorSdk.AFDOriginGroupsClient
	id     parse.FrontDoorOriginGroupId
}

var (
	frontDoorOriginGroupDeleteSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	frontDoorOriginGroupDeleteInProgress = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewFrontDoorOriginGroupDeletePoller(client *cdnFrontDoorSdk.AFDOriginGroupsClient, id parse.FrontDoorOriginGroupId) pollers.PollerType {
	return &frontDoorOriginGroupDeletePoller{
		client: client,
		id:     id,
	}
}

func (p *frontDoorOriginGroupDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id.ResourceGroup, p.id.ProfileName, p.id.OriginGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return &frontDoorOriginGroupDeleteSuccess, nil
		}

		return nil, fmt.Errorf("retrieving %s while waiting for deletion: %+v", p.id, err)
	}

	if props := resp.AFDOriginGroupProperties; props != nil {
		log.Printf("[DEBUG] AFD Origin Group %s still present while waiting for deletion (deploymentStatus=%q provisioningState=%q)", p.id, string(props.DeploymentStatus), string(props.ProvisioningState))
	}

	return &frontDoorOriginGroupDeleteInProgress, nil
}
