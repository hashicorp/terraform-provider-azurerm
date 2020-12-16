package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/response"
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

func resourceArmLogAnalyticsWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Read:   resourceArmLogAnalyticsWorkspaceRead,
		Update: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Delete: resourceArmLogAnalyticsWorkspaceDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LogAnalyticsWorkspaceID(id)
			return err
		}),

		SchemaVersion: 1,

		MigrateState: WorkspaceMigrateState,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"internet_ingestion_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"internet_query_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(operationalinsights.WorkspaceSkuNameEnumPerGB2018),
				ValidateFunc: validation.StringInSlice([]string{
					string(operationalinsights.WorkspaceSkuNameEnumFree),
					string(operationalinsights.WorkspaceSkuNameEnumPerGB2018),
					string(operationalinsights.WorkspaceSkuNameEnumPerNode),
					string(operationalinsights.WorkspaceSkuNameEnumPremium),
					string(operationalinsights.WorkspaceSkuNameEnumStandalone),
					string(operationalinsights.WorkspaceSkuNameEnumStandard),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, true),
				DiffSuppressFunc: logAnalyticsLinkedServiceSkuChangeCaseDifference,
			},

			"retention_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
			},

			"daily_quota_gb": {
				Type:             schema.TypeFloat,
				Optional:         true,
				Default:          -1.0,
				DiffSuppressFunc: dailyQuotaGbDiffSuppressFunc,
				ValidateFunc:     validation.FloatAtLeast(0),
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "this property has been removed from the API and will be removed in version 3.0 of the provider",
			},

			"primary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmLogAnalyticsWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Log Analytics Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	skuName := d.Get("sku").(string)
	sku := &operationalinsights.WorkspaceSku{
		Name: operationalinsights.WorkspaceSkuNameEnum(skuName),
	}

	// (@WodansSon) - If the workspace is connected to a cluster via the linked service resource
	// the workspace cannot be modified since the linked service changes the sku value within
	// the workspace
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
		if err == nil {
			if azSku := resp.Sku; azSku != nil {
				if strings.EqualFold(string(azSku.Name), "lacluster") {
					return fmt.Errorf("Log Analytics Workspace %q (Resource Group %q): cannot be modified while it is connected to a Log Analytics cluster", name, resourceGroup)
				}
			}
		}
	}

	internetIngestionEnabled := operationalinsights.Disabled
	if d.Get("internet_ingestion_enabled").(bool) {
		internetIngestionEnabled = operationalinsights.Enabled
	}
	internetQueryEnabled := operationalinsights.Disabled
	if d.Get("internet_query_enabled").(bool) {
		internetQueryEnabled = operationalinsights.Enabled
	}

	retentionInDays := int32(d.Get("retention_in_days").(int))

	t := d.Get("tags").(map[string]interface{})

	parameters := operationalinsights.Workspace{
		Name:     &name,
		Location: &location,
		Tags:     tags.Expand(t),
		WorkspaceProperties: &operationalinsights.WorkspaceProperties{
			Sku:                             sku,
			PublicNetworkAccessForIngestion: internetIngestionEnabled,
			PublicNetworkAccessForQuery:     internetQueryEnabled,
			RetentionInDays:                 &retentionInDays,
		},
	}

	dailyQuotaGb, ok := d.GetOk("daily_quota_gb")
	if ok && strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		return fmt.Errorf("`Free` tier SKU quota is not configurable and is hard set to 0.5GB")
	} else if !strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		parameters.WorkspaceProperties.WorkspaceCapping = &operationalinsights.WorkspaceCapping{
			DailyQuotaGb: utils.Float(dailyQuotaGb.(float64)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmLogAnalyticsWorkspaceRead(d, meta)
}

func resourceArmLogAnalyticsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	sharedKeysClient := meta.(*clients.Client).LogAnalytics.SharedKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics workspaces '%s': %+v", id.WorkspaceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("internet_ingestion_enabled", resp.PublicNetworkAccessForIngestion == operationalinsights.Enabled)
	d.Set("internet_query_enabled", resp.PublicNetworkAccessForQuery == operationalinsights.Enabled)

	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", "")
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)
	if resp.WorkspaceProperties != nil && resp.WorkspaceProperties.Sku != nil && strings.EqualFold(string(resp.WorkspaceProperties.Sku.Name), string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		// Special case for "Free" tier
		d.Set("daily_quota_gb", utils.Float(0.5))
	} else if workspaceCapping := resp.WorkspaceCapping; workspaceCapping != nil {
		d.Set("daily_quota_gb", resp.WorkspaceCapping.DailyQuotaGb)
	} else {
		d.Set("daily_quota_gb", utils.Float(-1))
	}

	sharedKeys, err := sharedKeysClient.GetSharedKeys(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", id.WorkspaceName, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogAnalyticsWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	force := false
	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, utils.Bool(force))
	if err != nil {
		return fmt.Errorf("issuing AzureRM delete request for Log Analytics Workspaces '%s': %+v", id.WorkspaceName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Log Analytics Worspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func dailyQuotaGbDiffSuppressFunc(_, _, _ string, d *schema.ResourceData) bool {
	// (@jackofallops) - 'free' is a legacy special case that is always set to 0.5GB
	if skuName := d.Get("sku").(string); strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		return true
	}

	return false
}

func logAnalyticsLinkedServiceSkuChangeCaseDifference(k, old, new string, d *schema.ResourceData) bool {
	// (@WodansSon) - This is needed because if you connect your workspace to a log analytics linked service resource it
	// will modify the value of your sku to "lacluster". We are currently in negotiations with the service team to
	// see if there is another way of doing this, for now this is the workaround
	if old == "lacluster" {
		old = new
	}

	return suppress.CaseDifference(k, old, new, d)
}
