// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &streamAnalyticsClusterDeletePoller{}

// streamAnalyticsClusterDeletePoller polls the Cluster itself until it's gone. The polling URL returned by the
// Delete request is nested under the Cluster, so it returns a 404 rather than a terminal status once the Cluster
// has been removed, which the long-running-operation poller treats as still in progress.
type streamAnalyticsClusterDeletePoller struct {
	client *clusters.ClustersClient
	id     clusters.ClusterId
}

func NewStreamAnalyticsClusterDeletePoller(client *clusters.ClustersClient, id clusters.ClusterId) pollers.PollerType {
	return &streamAnalyticsClusterDeletePoller{
		client: client,
		id:     id,
	}
}

func (p streamAnalyticsClusterDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusSucceeded,
			}, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
