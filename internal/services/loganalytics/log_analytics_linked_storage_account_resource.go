package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsLinkedStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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
			_, err := parse.LogAnalyticsLinkedStorageAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"data_source_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					strings.ToLower(string(operationalinsights.CustomLogs)),
					strings.ToLower(string(operationalinsights.AzureWatson)),
					strings.ToLower(string(operationalinsights.Query)),
					strings.ToLower(string(operationalinsights.Alerts)),
					// Value removed from enum in 2020-08-01, but effectively still works
					"ingestion",
				}, false),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
			},

			"storage_account_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}

func resourceLogAnalyticsLinkedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	id := parse.NewLogAnalyticsLinkedStorageAccountID(workspace.SubscriptionId, d.Get("resource_group_name").(string), workspace.WorkspaceName, d.Get("data_source_type").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, operationalinsights.DataSourceType(id.LinkedStorageAccountName))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_linked_storage_account", id.ID())
		}
	}

	parameters := operationalinsights.LinkedStorageAccountsResource{
		LinkedStorageAccountsProperties: &operationalinsights.LinkedStorageAccountsProperties{
			StorageAccountIds: utils.ExpandStringSlice(d.Get("storage_account_ids").(*pluginsdk.Set).List()),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, operationalinsights.DataSourceType(id.LinkedStorageAccountName), parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsLinkedStorageAccountRead(d, meta)
}

func resourceLogAnalyticsLinkedStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsLinkedStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	dataSourceType := operationalinsights.DataSourceType(id.LinkedStorageAccountName)
	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, dataSourceType)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics Linked Storage Account %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q): %+v", id.LinkedStorageAccountName, id.ResourceGroup, id.WorkspaceName, err)
	}

	d.Set("data_source_type", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_resource_id", parse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	if props := resp.LinkedStorageAccountsProperties; props != nil {
		d.Set("storage_account_ids", utils.FlattenStringSlice(props.StorageAccountIds))
	}

	return nil
}

func resourceLogAnalyticsLinkedStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsLinkedStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	dataSourceType := operationalinsights.DataSourceType(id.LinkedStorageAccountName)
	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, dataSourceType); err != nil {
		return fmt.Errorf("deleting Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q): %+v", id.LinkedStorageAccountName, id.ResourceGroup, id.WorkspaceName, err)
	}
	return nil
}
