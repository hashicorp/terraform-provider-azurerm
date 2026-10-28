// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_storage_insights

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsStorageInsightsCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	storageAccountId := d.Get("storage_account_id").(string)
	storageAccountKey := d.Get("storage_account_key").(string)

	workspace, err := storageinsights.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := storageinsights.NewStorageInsightConfigID(workspace.SubscriptionId, resourceGroup, workspace.WorkspaceName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.StorageInsightConfigsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for present of existing Log Analytics Storage Insights %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, id.WorkspaceName, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_log_analytics_storage_insights", id.ID())
			}
		}
	}

	parameters := storageinsights.StorageInsight{
		Properties: &storageinsights.StorageInsightProperties{
			StorageAccount: expandStorageInsightConfigStorageAccount(storageAccountId, storageAccountKey),
		},
	}

	if _, ok := d.GetOk("table_names"); ok {
		parameters.Properties.Tables = pluginsdk.ExpandStringSlice(d.Get("table_names").(*pluginsdk.Set).List())
	}

	if _, ok := d.GetOk("blob_container_names"); ok {
		parameters.Properties.Containers = pluginsdk.ExpandStringSlice(d.Get("blob_container_names").(*pluginsdk.Set).List())
	}

	if _, err := client.StorageInsightConfigsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsStorageInsightsRead(d, meta)
}

func expandStorageInsightConfigStorageAccount(id string, key string) storageinsights.StorageAccount {
	return storageinsights.StorageAccount{
		Id:  id,
		Key: key,
	}
}
