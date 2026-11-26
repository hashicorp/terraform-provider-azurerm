// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsWorkspaceRead(d *pluginsdk.ResourceData, meta any) error {
	sharedKeyClient := meta.(*clients.Client).LogAnalytics.SharedKeyWorkspacesClient
	client := meta.(*clients.Client).LogAnalytics.WorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceLogAnalyticsWorkspaceFlatten(ctx, sharedKeyClient, d, id, resp.Model, true)
}

const azureLogAnalyticsWorkspaceName = "azurerm_log_analytics_workspace"

func resourceLogAnalyticsWorkspaceCustomDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	// Since sku needs to be a force new if the sku changes we need to have this
	// custom diff here because when you link the workspace to a cluster the
	// cluster changes the sku to LACluster, so we need to ignore the change
	// if it is LACluster else invoke the ForceNew as before...

	if d.HasChange("sku") {
		old, new := d.GetChange("sku")
		log.Printf("[INFO] Log Analytics Workspace SKU: OLD: %q, NEW: %q", old, new)
		// If the old value is not LACluster(e.g. "") return ForceNew because they are
		// really changing the sku...
		changingFromLACluster := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumLACluster)) || strings.EqualFold(old.(string), "")
		// changing from capacity reservation to perGB does not force new when the last sku update date is more than 31-days ago.
		// to let users do the change, we do not force new in this case and let the API error out.
		changingFromCapacityReservationToPerGB := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) && strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight))
		changingFromPerGBToCapacityReservation := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight)) && strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumCapacityReservation))
		if !changingFromCapacityReservationToPerGB && !changingFromLACluster && !changingFromPerGBToCapacityReservation {
			d.ForceNew("sku")
		}

		// Creation or update workspace to `standard` or `premium` SKU is not allowed. Reference:
		// https://learn.microsoft.com/azure/azure-monitor/logs/cost-logs#standard-and-premium-pricing-tiers
		if strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumStandard)) ||
			strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumPremium)) {
			return fmt.Errorf("creation of log analytics workspaces with `Standard` or `Premium` SKUs is no longer supported by Azure - see https://learn.microsoft.com/azure/azure-monitor/logs/cost-logs#standard-and-premium-pricing-tiers")
		}
	}

	return nil
}
