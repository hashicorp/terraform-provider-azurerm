package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsStorageInsightConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsStorageInsightConfigCreateUpdate,
		Read:   resourceArmLogAnalyticsStorageInsightConfigRead,
		Update: resourceArmLogAnalyticsStorageInsightConfigCreateUpdate,
		Delete: resourceArmLogAnalyticsStorageInsightConfigDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LogAnalyticsStorageInsightConfigID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsStorageInsightConfigName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"blob_container_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"storage_account_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"table_names": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmLogAnalyticsStorageInsightConfigCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightConfigClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	workspace, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, workspace.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Log Analytics Storage Insight Config %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_storage_insight_config", *existing.ID)
		}
	}

	parameters := operationalinsights.StorageInsight{
		StorageInsightProperties: &operationalinsights.StorageInsightProperties{
			Containers:     utils.ExpandStringSlice(d.Get("blob_container_names").(*schema.Set).List()),
			StorageAccount: expandArmStorageInsightConfigStorageAccount(d.Get("storage_account_resource_id").([]interface{})),
			Tables:         utils.ExpandStringSlice(d.Get("table_names").(*schema.Set).List()),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.Name, name, parameters); err != nil {
		return fmt.Errorf("creating/updating Log Analytics Storage Insight Config %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspace.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Storage Insight Config %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Storage Insight Config %q (Resource Group %q / workspaceName %q) ID", name, resourceGroup, workspace.Name)
	}

	d.SetId(*resp.ID)
	return resourceArmLogAnalyticsStorageInsightConfigRead(d, meta)
}

func resourceArmLogAnalyticsStorageInsightConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightConfigClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsStorageInsightConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics Storage Insight Config %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Storage Insight Config %q (Resource Group %q / workspaceName %q): %+v", id.Name, id.ResourceGroup, id.WorkspaceName, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_name", id.WorkspaceName)
	d.Set("e_tag", resp.ETag)
	if props := resp.StorageInsightProperties; props != nil {
		d.Set("containers", utils.FlattenStringSlice(props.Containers))
		if err := d.Set("storage_account", flattenArmStorageInsightConfigStorageAccount(props.StorageAccount)); err != nil {
			return fmt.Errorf("setting `storage_account`: %+v", err)
		}
		d.Set("tables", utils.FlattenStringSlice(props.Tables))
	}
	d.Set("type", resp.Type)
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogAnalyticsStorageInsightConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightConfigClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsStorageInsightConfigID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting LogAnalytics Storage Insight Config %q (Resource Group %q / workspaceName %q): %+v", id.Name, id.ResourceGroup, id.WorkspaceName, err)
	}
	return nil
}

func expandArmStorageInsightConfigStorageAccount(input []interface{}) *operationalinsights.StorageAccount {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &operationalinsights.StorageAccount{
		Key: utils.String(v["key"].(string)),
	}
}

func flattenArmStorageInsightConfigStorageAccount(input *operationalinsights.StorageAccount) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var key string
	if input.Key != nil {
		key = *input.Key
	}
	return []interface{}{
		map[string]interface{}{
			"key": key,
		},
	}
}
