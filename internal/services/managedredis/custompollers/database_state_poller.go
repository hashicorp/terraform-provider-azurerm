// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dbStatePoller{}

type DatabasesClientInterface interface {
	Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error)
}

type dbStatePoller struct {
	client DatabasesClientInterface
	id     databases.DatabaseId
}

// Poll for db properties "resourceState" == "Running" as the LRO only polls for "provisioningState" and can
// complete prematurely
func NewDBStatePoller(client DatabasesClientInterface, id databases.DatabaseId) *dbStatePoller {
	return &dbStatePoller{
		client: client,
		id:     id,
	}
}

func (p dbStatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ResourceState == nil {
		return nil, fmt.Errorf("polling for %s: `resourceState` was empty", p.id)
	}

	resourceState := pointer.FromEnum(resp.Model.Properties.ResourceState)
	log.Printf("[DEBUG] DatabaseStatePoller for id: %s, resourceState: %s..", p.id, resourceState)

	if resourceState == string(databases.ResourceStateRunning) {
		return &pollingSuccess, nil
	}

	pendingStates := map[string]bool{
		string(databases.ResourceStateCreating):  true,
		string(databases.ResourceStateUpdating):  true,
		string(databases.ResourceStateEnabling):  true,
		string(databases.ResourceStateDeleting):  true,
		string(databases.ResourceStateDisabling): true,
		string(databases.ResourceStateMoving):    true,
		string(databases.ResourceStateScaling):   true,
	}
	if pendingStates[resourceState] {
		return &pollingInProgress, nil
	}

	return nil, fmt.Errorf("unexpected resource state for %s: %s", p.id, resourceState)
}
