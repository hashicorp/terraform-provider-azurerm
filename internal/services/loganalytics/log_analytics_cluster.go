// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func logAnalyticsClusterWaitForState(ctx context.Context, client *clusters.ClustersClient, clusterId clusters.ClusterId) (*pluginsdk.StateChangeConf, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context had no deadline")
	}
	return &pluginsdk.StateChangeConf{
		Pending:    []string{string(clusters.ClusterEntityStatusUpdating)},
		Target:     []string{string(clusters.ClusterEntityStatusSucceeded)},
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
		Refresh:    logAnalyticsClusterRefresh(ctx, client, clusterId),
	}, nil
}

func logAnalyticsClusterRefresh(ctx context.Context, client *clusters.ClustersClient, clusterId clusters.ClusterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] checking on state of Log Analytics Cluster %q", clusterId.ClusterName)

		resp, err := client.Get(ctx, clusterId)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of %q: %v", clusterId, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil {
			if resp.Model.Properties.ProvisioningState != nil {
				if *resp.Model.Properties.ProvisioningState != clusters.ClusterEntityStatusUpdating && *resp.Model.Properties.ProvisioningState != clusters.ClusterEntityStatusSucceeded {
					return nil, "nil", fmt.Errorf("%q unexpected Provisioning State encountered: %q", clusterId, string(*resp.Model.Properties.ProvisioningState))
				}
				return resp, string(*resp.Model.Properties.ProvisioningState), nil
			}
		}
		// I am not returning an error here as this might have just been a bad get
		return resp, "nil", nil
	}
}
