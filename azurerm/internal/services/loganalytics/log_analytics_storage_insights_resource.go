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
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsStorageInsights() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Read:   resourceLogAnalyticsStorageInsightsRead,
		Update: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Delete: resourceLogAnalyticsStorageInsightsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: logAnalyticsStorageInsightsImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsStorageInsightsName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"storage_account_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validate.IsBase64Encoded,
				),
			},

			"blob_container_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
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

func resourceLogAnalyticsStorageInsightsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	storageAccountId := d.Get("storage_account_id").(string)
	storageAccountKey := d.Get("storage_account_key").(string)

	workspaceId := d.Get("workspace_id").(string)
	id := parse.NewLogAnalyticsStorageInsightsId(resourceGroup, workspaceId, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, id.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Log Analytics Storage Insights %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, id.WorkspaceName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_storage_insights", *existing.ID)
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, id.WorkspaceName, name, parameters); err != nil {
		return fmt.Errorf("creating/updating Log Analytics Storage Insights %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, id.WorkspaceName, err)
	}

	d.SetId(id.ID(subscriptionId))
	return resourceLogAnalyticsStorageInsightsRead(d, meta)
}

func resourceLogAnalyticsStorageInsightsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsStorageInsightsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics Storage Insights %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Storage Insights %q (Resource Group %q / workspaceName %q): %+v", id.Name, id.ResourceGroup, id.WorkspaceName, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_id", id.WorkspaceID)

	if props := resp.StorageInsightProperties; props != nil {
		d.Set("blob_container_names", utils.FlattenStringSlice(props.Containers))
		storageAccountId := ""
		if props.StorageAccount != nil && props.StorageAccount.ID != nil {
			storageAccountId = *props.StorageAccount.ID
		}
		d.Set("storage_account_id", storageAccountId)
		d.Set("table_names", utils.FlattenStringSlice(props.Tables))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogAnalyticsStorageInsightsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.StorageInsightsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsStorageInsightsID(d.Id())
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
