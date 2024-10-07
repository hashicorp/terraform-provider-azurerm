// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsStorageInsights() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Read:   resourceLogAnalyticsStorageInsightsRead,
		Update: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Delete: resourceLogAnalyticsStorageInsightsDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := storageinsights.ParseStorageInsightConfigID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			if v, ok := d.GetOk("storage_account_key"); ok && v.(string) != "" {
				d.Set("storage_account_key", v)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: resourceLogAnalyticsStorageInsightsSchema(),
	}
}

func resourceLogAnalyticsStorageInsightsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
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
	id := storageinsights.NewStorageInsightConfigID(subscriptionId, resourceGroup, workspace.WorkspaceName, name)

	if d.IsNewResource() {
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

	parameters := storageinsights.StorageInsight{
		Properties: &storageinsights.StorageInsightProperties{
			StorageAccount: expandStorageInsightConfigStorageAccount(storageAccountId, storageAccountKey),
		},
	}

	if _, ok := d.GetOk("table_names"); ok {
		parameters.Properties.Tables = utils.ExpandStringSlice(d.Get("table_names").(*pluginsdk.Set).List())
	}

	if _, ok := d.GetOk("blob_container_names"); ok {
		parameters.Properties.Containers = utils.ExpandStringSlice(d.Get("blob_container_names").(*pluginsdk.Set).List())
	}

	if _, err := client.StorageInsightConfigsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsStorageInsightsRead(d, meta)
}

func resourceLogAnalyticsStorageInsightsRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("blob_container_names", utils.FlattenStringSlice(props.Containers))

			storageAccountIdStr := ""
			if props.StorageAccount.Id != "" {
				storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(props.StorageAccount.Id)
				if err != nil {
					return err
				}
				storageAccountIdStr = storageAccountId.ID()
			}
			d.Set("storage_account_id", storageAccountIdStr)

			d.Set("table_names", utils.FlattenStringSlice(props.Tables))
		}
	}

	return nil
}

func resourceLogAnalyticsStorageInsightsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storageinsights.ParseStorageInsightConfigID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.StorageInsightConfigsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandStorageInsightConfigStorageAccount(id string, key string) storageinsights.StorageAccount {
	return storageinsights.StorageAccount{
		Id:  id,
		Key: key,
	}
}

func resourceLogAnalyticsStorageInsightsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogAnalyticsStorageInsightsName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: storageinsights.ValidateWorkspaceID,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"storage_account_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: azValidate.Base64EncodedString,
		},

		"blob_container_names": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},

		"table_names": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
