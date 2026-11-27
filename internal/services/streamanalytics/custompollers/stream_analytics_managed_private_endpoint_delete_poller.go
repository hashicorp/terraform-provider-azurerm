// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &streamAnalyticsManagedPrivateEndpointDeletePoller{}

// streamAnalyticsManagedPrivateEndpointDeletePoller polls the Private Endpoint itself until it's gone, since the
// long-running-operation returned by the Delete request continues to report `InProgress` beyond the delete timeout.
type streamAnalyticsManagedPrivateEndpointDeletePoller struct {
	client *privateendpoints.PrivateEndpointsClient
	id     privateendpoints.PrivateEndpointId
}

func NewStreamAnalyticsManagedPrivateEndpointDeletePoller(client *privateendpoints.PrivateEndpointsClient, id privateendpoints.PrivateEndpointId) pollers.PollerType {
	return &streamAnalyticsManagedPrivateEndpointDeletePoller{
		client: client,
		id:     id,
	}
}

func (p streamAnalyticsManagedPrivateEndpointDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 10 * time.Second,
				Status:       pollers.PollingStatusSucceeded,
			}, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
