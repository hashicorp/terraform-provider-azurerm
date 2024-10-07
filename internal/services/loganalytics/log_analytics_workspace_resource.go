// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	sharedKeyWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
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
			_, err := workspaces.ParseWorkspaceID(id)
			return err
		}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
			1: migration.WorkspaceV1ToV2{},
			2: migration.WorkspaceV2ToV3{},
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

			"location": commonschema.Location(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"allow_resource_only_permissions": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"local_authentication_disabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"cmk_for_query_forced": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

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
					string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight),
					string(workspaces.WorkspaceSkuNameEnumPerNode),
					string(workspaces.WorkspaceSkuNameEnumPremium),
					string(workspaces.WorkspaceSkuNameEnumStandalone),
					string(workspaces.WorkspaceSkuNameEnumStandard),
					string(workspaces.WorkspaceSkuNameEnumCapacityReservation),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, false),
			},

			"reservation_capacity_in_gb_per_day": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(workspaces.CapacityReservationLevelOneHundred),
					int(workspaces.CapacityReservationLevelTwoHundred),
					int(workspaces.CapacityReservationLevelThreeHundred),
					int(workspaces.CapacityReservationLevelFourHundred),
					int(workspaces.CapacityReservationLevelFiveHundred),
					int(workspaces.CapacityReservationLevelOneThousand),
					int(workspaces.CapacityReservationLevelTwoThousand),
					int(workspaces.CapacityReservationLevelFiveThousand),
				}),
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

			"data_collection_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: datacollectionrules.ValidateDataCollectionRuleID,
			},

			"immediate_data_purge_on_30_days_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceLogAnalyticsWorkspaceCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	// Since sku needs to be a force new if the sku changes we need to have this
	// custom diff here because when you link the workspace to a cluster the
	// cluster changes the sku to LACluster, so we need to ignore the change
	// if it is LACluster else invoke the ForceNew as before...

	if d.HasChange("sku") {
		old, new := d.GetChange("sku")
		log.Printf("[INFO] Log Analytics Workspace SKU: OLD: %q, NEW: %q", old, new)
		// If the old value is not LACluster(e.g. "") return ForceNew because they are
		// really changing the sku...
		changingFromLACluster := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumLACluster)) || strings.EqualFold(old.(string), "")
		// changing from capacity reservation to perGB does not force new when the last sku update date is more than 31-days ago.
		// to let users do the change, we do not force new in this case and let the API error out.
		changingFromCapacityReservationToPerGB := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) && strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight))
		changingFromPerGBToCapacityReservation := strings.EqualFold(old.(string), string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight)) && strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumCapacityReservation))
		if !changingFromCapacityReservationToPerGB && !changingFromLACluster && !changingFromPerGBToCapacityReservation {
			d.ForceNew("sku")
		}
	}

	return nil
}

func resourceLogAnalyticsWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspaceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace creation.")

	var isLACluster bool
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := workspaces.NewWorkspaceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Log Analytics Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	skuName := d.Get("sku").(string)
	sku := &workspaces.WorkspaceSku{
		Name: workspaces.WorkspaceSkuNameEnum(skuName),
	}

	// (@WodansSon) - If the workspace is connected to a cluster via the linked service resource
	// the workspace SKU cannot be modified since the linked service owns the sku value within
	// the workspace once it is linked
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err == nil {
			if resp.Model != nil && resp.Model.Properties != nil {
				if azSku := resp.Model.Properties.Sku; azSku != nil {
					if strings.EqualFold(string(azSku.Name), string(workspaces.WorkspaceSkuNameEnumLACluster)) {
						isLACluster = true
						log.Printf("[INFO] Log Analytics Workspace %q (Resource Group %q): SKU is linked to Log Analytics cluster", name, resourceGroup)
					}
				}
			}
		}
	}

	internetIngestionEnabled := workspaces.PublicNetworkAccessTypeDisabled
	if d.Get("internet_ingestion_enabled").(bool) {
		internetIngestionEnabled = workspaces.PublicNetworkAccessTypeEnabled
	}
	internetQueryEnabled := workspaces.PublicNetworkAccessTypeDisabled
	if d.Get("internet_query_enabled").(bool) {
		internetQueryEnabled = workspaces.PublicNetworkAccessTypeEnabled
	}

	retentionInDays := int64(d.Get("retention_in_days").(int))

	t := d.Get("tags").(map[string]interface{})

	if isLACluster {
		sku.Name = workspaces.WorkspaceSkuNameEnumLACluster
	} else if skuName == "" {
		// Default value if sku is not defined
		sku.Name = workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight
	}

	allowResourceOnlyPermission := d.Get("allow_resource_only_permissions").(bool)
	disableLocalAuth := d.Get("local_authentication_disabled").(bool)

	parameters := workspaces.Workspace{
		Name:     &name,
		Location: location,
		Tags:     expandTags(t),
		Properties: &workspaces.WorkspaceProperties{
			Sku:                             sku,
			PublicNetworkAccessForIngestion: &internetIngestionEnabled,
			PublicNetworkAccessForQuery:     &internetQueryEnabled,
			RetentionInDays:                 &retentionInDays,
			Features: &workspaces.WorkspaceFeatures{
				EnableLogAccessUsingOnlyResourcePermissions: utils.Bool(allowResourceOnlyPermission),
				DisableLocalAuth: utils.Bool(disableLocalAuth),
			},
		},
	}

	// nolint : staticcheck
	if v, ok := d.GetOkExists("cmk_for_query_forced"); ok {
		parameters.Properties.ForceCmkForQuery = utils.Bool(v.(bool))
	}

	dailyQuotaGb, ok := d.GetOk("daily_quota_gb")
	if ok && strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) && (dailyQuotaGb != -1 && dailyQuotaGb != 0.5) {
		return fmt.Errorf("`Free` tier SKU quota is not configurable and is hard set to 0.5GB")
	} else if !strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) {
		parameters.Properties.WorkspaceCapping = &workspaces.WorkspaceCapping{
			DailyQuotaGb: utils.Float(dailyQuotaGb.(float64)),
		}
	}

	// The `ImmediatePurgeDataOn30Days` are not returned before it has been set
	// nolint : staticcheck
	if v, ok := d.GetOkExists("immediate_data_purge_on_30_days_enabled"); ok {
		parameters.Properties.Features.ImmediatePurgeDataOn30Days = utils.Bool(v.(bool))
	}

	propName := "reservation_capacity_in_gb_per_day"
	capacityReservationLevel, ok := d.GetOk(propName)
	if ok {
		if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			capacityReservationLevelValue := workspaces.CapacityReservationLevel(int64(capacityReservationLevel.(int)))
			parameters.Properties.Sku.CapacityReservationLevel = &capacityReservationLevelValue
		} else {
			return fmt.Errorf("`%s` can only be used with the `CapacityReservation` SKU", propName)
		}
	} else {
		if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			return fmt.Errorf("`%s` must be set when using the `CapacityReservation` SKU", propName)
		}
	}

	if v, ok := d.GetOk("identity"); ok {
		expanded, err := identity.ExpandSystemOrUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding identity: %+v", err)
		}
		parameters.Identity = expanded
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return err
	}

	// `data_collection_rule_id` also needs an additional update.
	// error message: Default dcr is not applicable on workspace creation, please provide it on update.
	if v, ok := d.GetOk("data_collection_rule_id"); ok {
		parameters.Properties.DefaultDataCollectionRuleResourceId = pointer.To(v.(string))
	}

	// `allow_resource_only_permissions` needs an additional update, tacked on https://github.com/Azure/azure-rest-api-specs/issues/21591
	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return err
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{strconv.FormatBool(!allowResourceOnlyPermission)},
		Target:     []string{strconv.FormatBool(allowResourceOnlyPermission)},
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				return resp, "error", fmt.Errorf("retiring %s: %+v", id, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Features != nil && resp.Model.Properties.Features.EnableLogAccessUsingOnlyResourcePermissions != nil {
				return resp, strconv.FormatBool(*resp.Model.Properties.Features.EnableLogAccessUsingOnlyResourcePermissions), nil
			}

			return resp, "false", fmt.Errorf("retiring %s: feature is nil", id)
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting on update for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsWorkspaceRead(d, meta)
}

func resourceLogAnalyticsWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	sharedKeyClient := meta.(*clients.Client).LogAnalytics.SharedKeyWorkspacesClient
	client := meta.(*clients.Client).LogAnalytics.WorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspaces '%s': %+v", id.WorkspaceName, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if model.Identity != nil {
			flattenIdentity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening identity: %+v", err)
			}
			d.Set("identity", flattenIdentity)
		}

		if props := model.Properties; props != nil {
			internetIngestionEnabled := true
			if props.PublicNetworkAccessForIngestion != nil {
				internetIngestionEnabled = *props.PublicNetworkAccessForIngestion == workspaces.PublicNetworkAccessTypeEnabled
			}
			d.Set("internet_ingestion_enabled", internetIngestionEnabled)

			internetQueryEnabled := true
			if props.PublicNetworkAccessForQuery != nil {
				internetQueryEnabled = *props.PublicNetworkAccessForQuery == workspaces.PublicNetworkAccessTypeEnabled
			}
			d.Set("internet_query_enabled", internetQueryEnabled)

			customerId := ""
			if props.CustomerId != nil {
				customerId = *props.CustomerId
			}
			d.Set("workspace_id", customerId)

			skuName := ""
			if props.Sku != nil {
				sku := *props.Sku
				for _, v := range workspaces.PossibleValuesForWorkspaceSkuNameEnum() {
					if strings.EqualFold(v, string(sku.Name)) {
						skuName = v
					}
				}
				if capacityReservationLevel := sku.CapacityReservationLevel; capacityReservationLevel != nil {
					d.Set("reservation_capacity_in_gb_per_day", int64(pointer.From(capacityReservationLevel)))
				}
			}
			d.Set("sku", skuName)

			forceCmkForQuery := false
			if props.ForceCmkForQuery != nil {
				forceCmkForQuery = *props.ForceCmkForQuery
			}
			d.Set("cmk_for_query_forced", forceCmkForQuery)

			var retentionInDays int64
			if props.RetentionInDays != nil {
				retentionInDays = *props.RetentionInDays
			}
			d.Set("retention_in_days", retentionInDays)

			switch {
			case strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)):
				// Special case for "Free" tier
				d.Set("daily_quota_gb", utils.Float(0.5))
			case props.WorkspaceCapping != nil && props.WorkspaceCapping.DailyQuotaGb != nil:
				d.Set("daily_quota_gb", props.WorkspaceCapping.DailyQuotaGb)
			default:
				d.Set("daily_quota_gb", utils.Float(-1))
			}

			allowResourceOnlyPermissions := true
			disableLocalAuth := false
			purgeDataOnThirtyDays := false
			if features := props.Features; features != nil {
				allowResourceOnlyPermissions = pointer.From(features.EnableLogAccessUsingOnlyResourcePermissions)
				disableLocalAuth = pointer.From(features.DisableLocalAuth)
				purgeDataOnThirtyDays = pointer.From(features.ImmediatePurgeDataOn30Days)
			}
			d.Set("allow_resource_only_permissions", allowResourceOnlyPermissions)
			d.Set("local_authentication_disabled", disableLocalAuth)
			d.Set("immediate_data_purge_on_30_days_enabled", purgeDataOnThirtyDays)

			defaultDataCollectionRuleResourceId := ""
			if props.DefaultDataCollectionRuleResourceId != nil {
				dataCollectionId, err := datacollectionrules.ParseDataCollectionRuleID(*props.DefaultDataCollectionRuleResourceId)
				if err != nil {
					return err
				}

				defaultDataCollectionRuleResourceId = dataCollectionId.ID()
			}
			d.Set("data_collection_rule_id", defaultDataCollectionRuleResourceId)

			sharedKeyId := sharedKeyWorkspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			sharedKeysResp, err := sharedKeyClient.SharedKeysGetSharedKeys(ctx, sharedKeyId)
			if err != nil {
				log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", id.WorkspaceName, err)
			} else {
				if sharedKeysModel := sharedKeysResp.Model; sharedKeysModel != nil {
					primarySharedKey := ""
					if sharedKeysModel.PrimarySharedKey != nil {
						primarySharedKey = *sharedKeysModel.PrimarySharedKey
					}
					d.Set("primary_shared_key", primarySharedKey)

					secondarySharedKey := ""
					if sharedKeysModel.SecondarySharedKey != nil {
						secondarySharedKey = *sharedKeysModel.SecondarySharedKey
					}
					d.Set("secondary_shared_key", secondarySharedKey)
				}
			}
		}

		d.Set("location", azure.NormalizeLocation(model.Location))

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceLogAnalyticsWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SharedKeyWorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	sharedKeyId := sharedKeyWorkspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return err
	}

	permanentlyDeleteOnDestroy := meta.(*clients.Client).Features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy
	err = client.DeleteThenPoll(ctx, sharedKeyId, sharedKeyWorkspaces.DeleteOperationOptions{Force: utils.Bool(permanentlyDeleteOnDestroy)})
	if err != nil {
		return fmt.Errorf("issuing AzureRM delete request for Log Analytics Workspaces '%s': %+v", id.WorkspaceName, err)
	}

	return nil
}

func dailyQuotaGbDiffSuppressFunc(_, _, _ string, d *pluginsdk.ResourceData) bool {
	// (@jackofallops) - 'free' is a legacy special case that is always set to 0.5GB
	if skuName := d.Get("sku").(string); strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) {
		return true
	}

	return false
}
