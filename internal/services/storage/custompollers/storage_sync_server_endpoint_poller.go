// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2022-09-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &storageSyncServerEndpointPoller{}

// This resource has 2 potential 'ready' statuses
//
// we shouldn't require that the file sync client has been configured correctly in order to declare success
//
// serverendpointresource.ServerProvisioningStatusReadySyncFunctional = "Ready_SyncFunctional"
// serverendpointresource.ServerProvisioningStatusReadySyncNotFunctional = "Ready_SyncNotFunctional"
type storageSyncServerEndpointPoller struct {
	client *serverendpointresource.ServerEndpointResourceClient
	id     serverendpointresource.ServerEndpointId
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatus("Ready_"),
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func NewStorageSyncServerEndpointPoller(client *serverendpointresource.ServerEndpointResourceClient, id serverendpointresource.ServerEndpointId) *storageSyncServerEndpointPoller {
	return &storageSyncServerEndpointPoller{
		client: client,
		id:     id,
	}
}

func (p storageSyncServerEndpointPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.ServerEndpointsGet(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model != nil {
		if provisioningStatus := resp.Model.Properties.ProvisioningState; provisioningStatus != nil {
			if !strings.Contains(strings.ToLower((*provisioningStatus)), strings.ToLower(string(pollingSuccess.Status))) {
				return &pollingInProgress, nil
			}
		}
	}

	return &pollingSuccess, nil
}
