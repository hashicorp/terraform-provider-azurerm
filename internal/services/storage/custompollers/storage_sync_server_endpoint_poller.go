// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &storageSyncServerEndpointPoller{}

type storageSyncServerEndpointPoller struct {
	client *serverendpointresource.ServerEndpointResourceClient
	id     serverendpointresource.ServerEndpointId
}

// The `ServerEndpointsCreateThenPoll` and `ServerEndpointsUpdateThenPoll` methods do not properly await, so a custom poller is required
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

	provisioningState := ""
	if model := resp.Model; model != nil && model.Properties != nil {
		provisioningState = pointer.From(model.Properties.ProvisioningState)
	}

	if strings.EqualFold(provisioningState, "Succeeded") {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	if strings.EqualFold(provisioningState, "runServerJob") {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return nil, pollers.PollingFailedError{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		Message: fmt.Sprintf("unexpected provisioningState %q", provisioningState),
	}
}
