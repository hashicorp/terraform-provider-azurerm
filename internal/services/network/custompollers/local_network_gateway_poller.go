package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &localNetworkGatewayPoller{}

type localNetworkGatewayPoller struct {
	client *localnetworkgateways.LocalNetworkGatewaysClient
	id     localnetworkgateways.LocalNetworkGatewayId
}

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

	pollingResult := pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}

	if resp.Model != nil {
		if pointer.From(resp.Model.Properties.ProvisioningState) == localnetworkgateways.ProvisioningStateSucceeded {
			pollingResult.Status = pollers.PollingStatusSucceeded
			return &pollingResult, nil
		}
	}

	return &pollingResult, nil
}
