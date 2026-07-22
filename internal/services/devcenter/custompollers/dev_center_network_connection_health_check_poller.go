// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &devCenterNetworkConnectionHealthCheckPoller{}

type devCenterNetworkConnectionHealthCheckPoller struct {
	client *networkconnections.NetworkConnectionsClient
	id     networkconnections.NetworkConnectionId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewDevCenterNetworkConnectionHealthCheckPoller(client *networkconnections.NetworkConnectionsClient, id networkconnections.NetworkConnectionId) *devCenterNetworkConnectionHealthCheckPoller {
	return &devCenterNetworkConnectionHealthCheckPoller{
		client: client,
		id:     id,
	}
}

func (p devCenterNetworkConnectionHealthCheckPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HealthCheckStatus == nil {
		return nil, fmt.Errorf("polling for %s: `healthCheckStatus` was empty", p.id)
	}

	if string(*resp.Model.Properties.HealthCheckStatus) == string(networkconnections.HealthCheckStatusPassed) {
		return &pollingSuccess, nil
	}

	return &pollingInProgress, nil
}
