// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	loganalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsLinkedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := loganalyticsClient.From(meta).LinkedStorageAccountClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	workspaceId := d.Get("workspace_id").(string)
	workspace, err := linkedstorageaccounts.ParseWorkspaceID(workspaceId)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	id := linkedstorageaccounts.NewDataSourceTypeID(workspace.SubscriptionId, d.Get("resource_group_name").(string), workspace.WorkspaceName, linkedstorageaccounts.DataSourceType(d.Get("data_source_type").(string)))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_log_analytics_linked_storage_account", id.ID())
			}
		}
	}

	parameters := linkedstorageaccounts.LinkedStorageAccountsResource{
		Properties: linkedstorageaccounts.LinkedStorageAccountsProperties{
			StorageAccountIds: pluginsdk.ExpandStringSlice(d.Get("storage_account_ids").(*pluginsdk.Set).List()),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceLogAnalyticsLinkedStorageAccountRead(d, meta)
}
