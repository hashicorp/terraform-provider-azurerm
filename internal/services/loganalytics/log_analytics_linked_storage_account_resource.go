// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsLinkedStorageAccount() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceLogAnalyticsLinkedStorageAccountCreateUpdate,
		Read:   resourceLogAnalyticsLinkedStorageAccountRead,
		Update: resourceLogAnalyticsLinkedStorageAccountCreateUpdate,
		Delete: resourceLogAnalyticsLinkedStorageAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := linkedstorageaccounts.ParseDataSourceTypeID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LinkedStorageAccountV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"data_source_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(linkedstorageaccounts.DataSourceTypeCustomLogs),
					string(linkedstorageaccounts.DataSourceTypeAzureWatson),
					string(linkedstorageaccounts.DataSourceTypeQuery),
					string(linkedstorageaccounts.DataSourceTypeAlerts),
					string(linkedstorageaccounts.DataSourceTypeIngestion),
				}, !features.FourPointOhBeta()),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			// TODO: rename to `workspace_id` in 4.0
			"workspace_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: linkedstorageaccounts.ValidateWorkspaceID,
			},

			"storage_account_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: commonids.ValidateStorageAccountID,
				},
			},
		},
	}

	if !features.FourPointOh() {
		resource.Schema["data_source_type"].DiffSuppressFunc = suppress.CaseDifference
	}

	return resource
}

func resourceLogAnalyticsLinkedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := linkedstorageaccounts.ParseWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	id := linkedstorageaccounts.NewDataSourceTypeID(workspace.SubscriptionId, d.Get("resource_group_name").(string), workspace.WorkspaceName, linkedstorageaccounts.DataSourceType(d.Get("data_source_type").(string)))

	if d.IsNewResource() {
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

	parameters := linkedstorageaccounts.LinkedStorageAccountsResource{
		Properties: linkedstorageaccounts.LinkedStorageAccountsProperties{
			StorageAccountIds: utils.ExpandStringSlice(d.Get("storage_account_ids").(*pluginsdk.Set).List()),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsLinkedStorageAccountRead(d, meta)
}

func resourceLogAnalyticsLinkedStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedstorageaccounts.ParseDataSourceTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Log Analytics Linked Storage Account %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_resource_id", linkedstorageaccounts.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())
	d.Set("data_source_type", string(id.DataSourceType))

	if model := resp.Model; model != nil {
		props := model.Properties
		var storageAccountIds []string
		if props.StorageAccountIds != nil {
			storageAccountIds = *props.StorageAccountIds
		}
		d.Set("storage_account_ids", storageAccountIds)
	}

	return nil
}

func resourceLogAnalyticsLinkedStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedstorageaccounts.ParseDataSourceTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
