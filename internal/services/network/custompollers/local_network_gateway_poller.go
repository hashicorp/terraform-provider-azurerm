// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &localNetworkGatewayPoller{}

type localNetworkGatewayPoller struct {
	client *localnetworkgateways.LocalNetworkGatewaysClient
	id     localnetworkgateways.LocalNetworkGatewayId
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

func NewLocalNetworkGatewayPoller(client *localnetworkgateways.LocalNetworkGatewaysClient, id localnetworkgateways.LocalNetworkGatewayId) *localNetworkGatewayPoller {
	return &localNetworkGatewayPoller{
		client: client,
		id:     id,
	}
}

func (p localNetworkGatewayPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model != nil {
		if provisioningStatus := resp.Model.Properties.ProvisioningState; provisioningStatus != nil {
			if !strings.EqualFold(string(*provisioningStatus), string(pollingSuccess.Status)) {
				return &pollingInProgress, nil
			}
		}
	}

	return &pollingSuccess, nil
}
