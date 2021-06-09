package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func logAnalyticsLinkedServiceDeleteWaitForState(ctx context.Context, meta interface{}, timeout time.Duration, resourceGroup string, workspaceName string, serviceType string) *pluginsdk.StateChangeConf {
	return &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh:    logAnalyticsLinkedServiceRefresh(ctx, meta, resourceGroup, workspaceName, serviceType),
	}
}

func logAnalyticsLinkedServiceRefresh(ctx context.Context, meta interface{}, resourceGroup string, workspaceName string, serviceType string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient

		log.Printf("[INFO] checking on state of Log Analytics Linked Service '%s/%s' (Resource Group %q)", workspaceName, serviceType, resourceGroup)

		resp, err := client.Get(ctx, resourceGroup, workspaceName, serviceType)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of Log Analytics Linked Service '%s/%s' (Resource Group %q)", workspaceName, serviceType, resourceGroup)
		}

		// (@WodansSon) - The service returns status code 200 even if the resource does not exist
		// instead it returns an empty slice...
		if props := resp.LinkedServiceProperties; props == nil {
			return resp, "Deleted", nil
		}

		return resp, "Deleting", nil
	}
}
