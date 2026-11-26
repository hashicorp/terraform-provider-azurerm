// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_storage_insights

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsStorageInsightsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storageinsights.ParseStorageInsightConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.StorageInsightConfigsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Log Analytics Storage Insights %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.StorageInsightConfigName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_id", storageinsights.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("blob_container_names", pluginsdk.FlattenSlice(props.Containers))

			storageAccountIdStr := ""
			if props.StorageAccount.Id != "" {
				storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(props.StorageAccount.Id)
				if err != nil {
					return err
				}
				storageAccountIdStr = storageAccountId.ID()
			}
			d.Set("storage_account_id", storageAccountIdStr)

			d.Set("table_names", pluginsdk.FlattenSlice(props.Tables))
		}
	}

	return nil
}
