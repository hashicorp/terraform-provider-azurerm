package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsWorkspaceCreateUpdate,
		Read:   resourceLogAnalyticsWorkspaceRead,
		Update: resourceLogAnalyticsWorkspaceCreateUpdate,
		Delete: resourceLogAnalyticsWorkspaceDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceLogAnalyticsWorkspaceCustomDiff),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogAnalyticsWorkspaceID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
			1: migration.WorkspaceV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"internet_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"internet_query_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// TODO 4.0: Clean up lacluster "workaround" to make it more readable and easier to understand. (@WodansSon already has the code written for the clean up)
			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(operationalinsights.WorkspaceSkuNameEnumFree),
					string(operationalinsights.WorkspaceSkuNameEnumPerGB2018),
					string(operationalinsights.WorkspaceSkuNameEnumPerNode),
					string(operationalinsights.WorkspaceSkuNameEnumPremium),
					string(operationalinsights.WorkspaceSkuNameEnumStandalone),
					string(operationalinsights.WorkspaceSkuNameEnumStandard),
					string(operationalinsights.WorkspaceSkuNameEnumCapacityReservation),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, false),
			},

			"reservation_capacity_in_gb_per_day": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntBetween(100, 5000), validation.IntDivisibleBy(100)),
			},

			"retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
			},

			"daily_quota_gb": {
				Type:             pluginsdk.TypeFloat,
				Optional:         true,
				Default:          -1.0,
				DiffSuppressFunc: dailyQuotaGbDiffSuppressFunc,
				ValidateFunc:     validation.FloatAtLeast(-1.0),
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceLogAnalyticsWorkspaceCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	// Since sku needs to be a force new if the sku changes we need to have this
	// custom diff here because when you link the workspace to a cluster the
	// cluster changes the sku to LACluster, so we need to ignore the change
	// if it is LACluster else invoke the ForceNew as before...
	//
	// NOTE: Since LACluster is not in our enum the value is returned as ""
	if d.HasChange("sku") {
		old, new := d.GetChange("sku")
		log.Printf("[INFO] Log Analytics Workspace SKU: OLD: %q, NEW: %q", old, new)
		// If the old value is not LACluster(e.g. "") return ForceNew because they are
		// really changing the sku...
		if !strings.EqualFold(old.(string), "") {
			d.ForceNew("sku")
		}
	}

	return nil
}

func resourceLogAnalyticsWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace creation.")

	var isLACluster bool
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Log Analytics Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	skuName := d.Get("sku").(string)
	sku := &operationalinsights.WorkspaceSku{
		Name: operationalinsights.WorkspaceSkuNameEnum(skuName),
	}

	// (@WodansSon) - If the workspace is connected to a cluster via the linked service resource
	// the workspace SKU cannot be modified since the linked service owns the sku value within
	// the workspace once it is linked
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
		if err == nil {
			if azSku := resp.Sku; azSku != nil {
				if strings.EqualFold(string(azSku.Name), "lacluster") {
					isLACluster = true
					log.Printf("[INFO] Log Analytics Workspace %q (Resource Group %q): SKU is linked to Log Analytics cluster", name, resourceGroup)
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

	if isLACluster {
		sku.Name = "lacluster"
	} else if skuName == "" {
		// Default value if sku is not defined
		sku.Name = operationalinsights.WorkspaceSkuNameEnumPerGB2018
	}

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
	if ok && strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) && (dailyQuotaGb != -1 && dailyQuotaGb != 0.5) {
		return fmt.Errorf("`Free` tier SKU quota is not configurable and is hard set to 0.5GB")
	} else if !strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		parameters.WorkspaceProperties.WorkspaceCapping = &operationalinsights.WorkspaceCapping{
			DailyQuotaGb: utils.Float(dailyQuotaGb.(float64)),
		}
	}

	propName := "reservation_capacity_in_gb_per_day"
	capacityReservationLevel, ok := d.GetOk(propName)
	if ok {
		if strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumCapacityReservation)) {
			parameters.WorkspaceProperties.Sku.CapacityReservationLevel = utils.Int32((int32(capacityReservationLevel.(int))))
		} else {
			return fmt.Errorf("`%s` can only be used with the `CapacityReservation` SKU", propName)
		}
	} else {
		if strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumCapacityReservation)) {
			return fmt.Errorf("`%s` must be set when using the `CapacityReservation` SKU", propName)
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

	return resourceLogAnalyticsWorkspaceRead(d, meta)
}

func resourceLogAnalyticsWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspaces '%s': %+v", id.WorkspaceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("internet_ingestion_enabled", resp.PublicNetworkAccessForIngestion == operationalinsights.Enabled)
	d.Set("internet_query_enabled", resp.PublicNetworkAccessForQuery == operationalinsights.Enabled)

	d.Set("workspace_id", resp.CustomerID)
	skuName := ""
	if sku := resp.Sku; sku != nil {
		for _, v := range operationalinsights.PossibleSkuNameEnumValues() {
			if strings.EqualFold(string(v), string(sku.Name)) {
				skuName = string(v)
			}
		}

		if capacityReservationLevel := sku.CapacityReservationLevel; capacityReservationLevel != nil {
			d.Set("reservation_capacity_in_gb_per_day", capacityReservationLevel)
		}
	}
	d.Set("sku", skuName)

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

func resourceLogAnalyticsWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsWorkspaceID(d.Id())
	if err != nil {
		return err
	}
	PermanentlyDeleteOnDestroy := meta.(*clients.Client).Features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy
	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, utils.Bool(PermanentlyDeleteOnDestroy))
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

func dailyQuotaGbDiffSuppressFunc(_, _, _ string, d *pluginsdk.ResourceData) bool {
	// (@jackofallops) - 'free' is a legacy special case that is always set to 0.5GB
	if skuName := d.Get("sku").(string); strings.EqualFold(skuName, string(operationalinsights.WorkspaceSkuNameEnumFree)) {
		return true
	}

	return false
}
