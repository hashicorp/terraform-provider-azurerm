package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
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

			// must ignore case since API lowercases all returned data
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"workspace_resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"storage_account_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validate.IsBase64Encoded,
			},

			"blob_container_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	storageAccountId := d.Get("storage_account_resource_id").(string)
	storageAccountKey := d.Get("storage_account_key").(string)
	if len(strings.TrimSpace(storageAccountKey)) < 1 {
		return fmt.Errorf("The argument 'storage_account_key' is required, but no definition was found.")
	}

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
			StorageAccount: expandArmStorageInsightConfigStorageAccount(storageAccountId, storageAccountKey),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("table_names"); ok {
		parameters.StorageInsightProperties.Tables = utils.ExpandStringSlice(d.Get("table_names").(*schema.Set).List())
	}

	if _, ok := d.GetOk("blob_container_names"); ok {
		parameters.StorageInsightProperties.Containers = utils.ExpandStringSlice(d.Get("blob_container_names").(*schema.Set).List())
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

	// Need to pull this from the config since the API does not return this value
	storageAccountKey := d.Get("storage_account_key").(string)

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
	d.Set("workspace_resource_id", id.WorkspaceID)

	if props := resp.StorageInsightProperties; props != nil {
		d.Set("blob_container_names", utils.FlattenStringSlice(props.Containers))
		d.Set("storage_account_resource_id", props.StorageAccount.ID)
		d.Set("storage_account_key", storageAccountKey)
		d.Set("table_names", utils.FlattenStringSlice(props.Tables))
	}

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

func expandArmStorageInsightConfigStorageAccount(id string, key string) *operationalinsights.StorageAccount {

	return &operationalinsights.StorageAccount{
		ID:  utils.String(id),
		Key: utils.String(key),
	}
}
