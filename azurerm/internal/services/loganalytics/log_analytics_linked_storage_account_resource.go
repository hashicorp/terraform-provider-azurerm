package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsLinkedStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsLinkedStorageAccountCreateUpdate,
		Read:   resourceArmLogAnalyticsLinkedStorageAccountRead,
		Update: resourceArmLogAnalyticsLinkedStorageAccountCreateUpdate,
		Delete: resourceArmLogAnalyticsLinkedStorageAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LogAnalyticsLinkedStorageAccountID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"data_source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					strings.ToLower(string(operationalinsights.CustomLogs)),
					strings.ToLower(string(operationalinsights.AzureWatson)),
					strings.ToLower(string(operationalinsights.Query)),
					strings.ToLower(string(operationalinsights.Alerts)),
					// Value removed from enum in 2020-08-01, but effectively still works
					"Ingestion",
				}, false),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
			},

			"storage_account_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}

func resourceArmLogAnalyticsLinkedStorageAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedStorageAccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataSourceType := operationalinsights.DataSourceType(d.Get("data_source_type").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	workspace, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, workspace.WorkspaceName, dataSourceType)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q): %+v", dataSourceType, resourceGroup, workspace.WorkspaceName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_linked_storage_account", *existing.ID)
		}
	}

	parameters := operationalinsights.LinkedStorageAccountsResource{
		LinkedStorageAccountsProperties: &operationalinsights.LinkedStorageAccountsProperties{
			StorageAccountIds: utils.ExpandStringSlice(d.Get("storage_account_ids").(*schema.Set).List()),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.WorkspaceName, dataSourceType, parameters); err != nil {
		return fmt.Errorf("creating/updating Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q): %+v", dataSourceType, resourceGroup, workspace.WorkspaceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspace.WorkspaceName, dataSourceType)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q): %+v", dataSourceType, resourceGroup, workspace.WorkspaceName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Linked Storage Account %q (Resource Group %q / workspaceName %q) ID", dataSourceType, resourceGroup, workspace.WorkspaceName)
	}

	d.SetId(*resp.ID)
	return resourceArmLogAnalyticsLinkedStorageAccountRead(d, meta)
}

func resourceArmLogAnalyticsLinkedStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("workspace_resource_id", parse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID(""))
	if props := resp.LinkedStorageAccountsProperties; props != nil {
		d.Set("storage_account_ids", utils.FlattenStringSlice(props.StorageAccountIds))
	}

	return nil
}

func resourceArmLogAnalyticsLinkedStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
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
