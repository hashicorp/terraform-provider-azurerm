// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &clusterStatePoller{}

type RedisEnterpriseClientInterface interface {
	Get(ctx context.Context, id redisenterprise.RedisEnterpriseId) (redisenterprise.GetOperationResponse, error)
}

type clusterStatePoller struct {
	client RedisEnterpriseClientInterface
	id     redisenterprise.RedisEnterpriseId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

// Poll for cluster properties "resourceState" == "Running" as the LRO only polls for "provisioningState" and can
// complete prematurely
func NewClusterStatePoller(client RedisEnterpriseClientInterface, id redisenterprise.RedisEnterpriseId) *clusterStatePoller {
	return &clusterStatePoller{
		client: client,
		id:     id,
	}
}

func (p clusterStatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ResourceState == nil {
		return nil, fmt.Errorf("polling for %s: `resourceState` was empty", p.id)
	}

	resourceState := pointer.FromEnum(resp.Model.Properties.ResourceState)
	log.Printf("[DEBUG] ClusterStatePoller for id: %s, resourceState: %s..", p.id, resourceState)

	if resourceState == "Running" {
		return &pollingSuccess, nil
	}

	pendingStates := map[string]bool{"Creating": true, "Updating": true, "Enabling": true, "Deleting": true, "Disabling": true}
	if pendingStates[resourceState] {
		return &pollingInProgress, nil
	}

	return nil, fmt.Errorf("unexpected resource state for %s: %s", p.id, resourceState)
}
