// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &EnableAzureMonitorPoller{}

// EnableAzureMonitorPoller polls until the Azure Monitor for the specified HDInsight Cluster has been enabled
// This works around an issue outlined in  https://github.com/hashicorp/go-azure-sdk/issues/518 where the API
// is a LRO which doesn't use `provisioningState`.
type EnableAzureMonitorPoller struct {
	client    *extensions.ExtensionsClient
	clusterId commonids.HDInsightClusterId
}

func NewEnableAzureMonitorPoller(client *extensions.ExtensionsClient, clusterId commonids.HDInsightClusterId) *EnableAzureMonitorPoller {
	return &EnableAzureMonitorPoller{
		client:    client,
		clusterId: clusterId,
	}
}

func (p *EnableAzureMonitorPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.GetAzureMonitorStatus(ctx, p.clusterId)
	if err != nil {
		return nil, fmt.Errorf("retrieving Azure Monitor Status for %s: %+v", p.clusterId, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving Azure Monitor Status for %s: `model` was nil", p.clusterId)
	}
	if resp.Model.ClusterMonitoringEnabled == nil {
		return nil, fmt.Errorf("retrieving Azure Monitor Status for %s: `model.ClusterMonitoringEnabled` was nil", p.clusterId)
	}

	status := pollers.PollingStatusInProgress
	if *resp.Model.ClusterMonitoringEnabled {
		status = pollers.PollingStatusSucceeded
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			OData:    resp.OData,
			Response: resp.HttpResponse,
		},
		PollInterval: 10 * time.Second,
		Status:       status,
	}, nil
}
