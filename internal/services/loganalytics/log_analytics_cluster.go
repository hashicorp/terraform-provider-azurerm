package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func logAnalyticsClusterWaitForState(ctx context.Context, meta interface{}, resourceGroup string, clusterName string) (*pluginsdk.StateChangeConf, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context had no deadline")
	}
	return &pluginsdk.StateChangeConf{
		Pending:    []string{string(operationalinsights.Updating)},
		Target:     []string{string(operationalinsights.Succeeded)},
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
		Refresh:    logAnalyticsClusterRefresh(ctx, meta, resourceGroup, clusterName),
	}, nil
}

func logAnalyticsClusterRefresh(ctx context.Context, meta interface{}, resourceGroup string, clusterName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*clients.Client).LogAnalytics.ClusterClient

		log.Printf("[INFO] checking on state of Log Analytics Cluster %q", clusterName)

		resp, err := client.Get(ctx, resourceGroup, clusterName)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of Log Analytics Cluster %q (Resource Group %q): %v", clusterName, resourceGroup, err)
		}

		if resp.ClusterProperties != nil {
			if resp.ClusterProperties.ProvisioningState != operationalinsights.Updating && resp.ClusterProperties.ProvisioningState != operationalinsights.Succeeded {
				return nil, "nil", fmt.Errorf("log analytics Cluster %q (Resource Group %q) unexpected Provisioning State encountered: %q", clusterName, resourceGroup, string(resp.ClusterProperties.ProvisioningState))
			}

			return resp, string(resp.ClusterProperties.ProvisioningState), nil
		}

		// I am not returning an error here as this might have just been a bad get
		return resp, "nil", nil
	}
}
