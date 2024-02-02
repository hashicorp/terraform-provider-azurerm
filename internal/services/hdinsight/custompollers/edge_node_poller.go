// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &EdgeNodePoller{}

// EdgeNodePoller polls until the Edge Nodes have finished provisioning
type EdgeNodePoller struct {
	client    *clusters.ClustersClient
	clusterId commonids.HDInsightClusterId
}

func NewEdgeNodePoller(client *clusters.ClustersClient, clusterId commonids.HDInsightClusterId) *EdgeNodePoller {
	return &EdgeNodePoller{
		client:    client,
		clusterId: clusterId,
	}
}

func (p *EdgeNodePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.clusterId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.clusterId, err)
	}

	clusterState := ""
	if model := resp.Model; model != nil && model.Properties != nil {
		clusterState = pointer.From(model.Properties.ClusterState)
	}

	if strings.EqualFold(clusterState, "Running") {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	if strings.EqualFold(clusterState, "AzureVMConfiguration") || strings.EqualFold(clusterState, "Accepted") || strings.EqualFold(clusterState, "HdInsightConfiguration") {
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
		Message: fmt.Sprintf("unexpected clusterState %q", clusterState),
	}
}
