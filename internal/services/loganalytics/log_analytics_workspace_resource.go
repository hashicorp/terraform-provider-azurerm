// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	sharedKeyWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity

const azureLogAnalyticsWorkspaceName = "azurerm_log_analytics_workspace"

func resourceLogAnalyticsWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsWorkspaceCreate,
		Read:   resourceLogAnalyticsWorkspaceRead,
		Update: resourceLogAnalyticsWorkspaceUpdate,
		Delete: resourceLogAnalyticsWorkspaceDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceLogAnalyticsWorkspaceCustomDiff),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&workspaces.WorkspaceId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&workspaces.WorkspaceId{}),

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

			"resource_group_name": commonschema.ResourceGroupName(),

			"allow_resource_only_permissions": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"cmk_for_query_forced": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"internet_ingestion_access_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(workspaces.PublicNetworkAccessTypeEnabled),
				ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForPublicNetworkAccessType(), false),
			},

			"internet_query_access_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(workspaces.PublicNetworkAccessTypeEnabled),
				ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForPublicNetworkAccessType(), false),
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true, // azignore:AZS007 - pre-existing violation
				ValidateFunc: validation.StringInSlice([]string{
					// creation with `Premium` and `Standard` is not allowed on API side, keep these values for existing resources.
					string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight),
					string(workspaces.WorkspaceSkuNameEnumPerNode),
					string(workspaces.WorkspaceSkuNameEnumPremium),
					string(workspaces.WorkspaceSkuNameEnumStandalone),
					string(workspaces.WorkspaceSkuNameEnumStandard),
					string(workspaces.WorkspaceSkuNameEnumCapacityReservation),
					string(workspaces.WorkspaceSkuNameEnumLACluster),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, false),
			},

			"reservation_capacity_in_gb_per_day": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					100,
					200,
					300,
					400,
					500,
					1000,
					2000,
					5000,
					10000,
					25000,
					50000,
				}),
			},

			"retention_in_days": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				// Note: O+C because Azure assigns a default retention period when not specified
				Computed:     true,
				ValidateFunc: validation.IntBetween(30, 730),
			},

			"daily_quota_gb": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      -1.0,
				ValidateFunc: validation.FloatAtLeast(-1.0),
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

func resourceLogAnalyticsWorkspaceCustomDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
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

		// Creation or update workspace to `standard` or `premium` SKU is not allowed. Reference:
		// https://learn.microsoft.com/en-us/azure/azure-monitor/logs/cost-logs#standard-and-premium-pricing-tiers
		if strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumStandard)) ||
			strings.EqualFold(new.(string), string(workspaces.WorkspaceSkuNameEnumPremium)) {
			return fmt.Errorf("creation of log analytics workspaces with `Standard` or `Premium` SKUs is no longer supported by Azure - see https://learn.microsoft.com/en-us/azure/azure-monitor/logs/cost-logs#standard-and-premium-pricing-tiers")
		}
	}

	return nil
}

func resourceLogAnalyticsWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspaceClient
	deletedWorkspaceClient := meta.(*clients.Client).LogAnalytics.DeletedWorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var isLACluster bool

	name := d.Get("name").(string)
	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", id.ID())
		}
	}

	deleted, err := deletedWorkspaceClient.List(ctx, commonids.NewSubscriptionID(id.SubscriptionId))
	if err != nil {
		return fmt.Errorf("listing deleted Log Analytics Workspaces: %+v", err)
	}

	if model := deleted.Model; model != nil {
		for _, v := range *model {
			if props := v.Properties; props != nil && props.Sku != nil {
				if pointer.From(v.Name) == name && string(props.Sku.Name) == string(workspaces.WorkspaceSkuNameEnumLACluster) {
					isLACluster = true
				}
			}
		}
	}

	skuName := d.Get("sku").(string)
	sku := &workspaces.WorkspaceSku{
		Name: workspaces.WorkspaceSkuNameEnum(skuName),
	}

	if !isLACluster && strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumLACluster)) {
		return fmt.Errorf("`sku` cannot be set to `LACluster` during creation unless the workspace is in a soft-deleted state while linked to a Log Analytics Cluster")
	}

	if isLACluster {
		sku.Name = workspaces.WorkspaceSkuNameEnumLACluster
	} else if skuName == "" {
		// Default value if sku is not defined
		sku.Name = workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight
	}

	allowResourceOnlyPermission := d.Get("allow_resource_only_permissions").(bool)

	parameters := workspaces.Workspace{
		Name:     &name,
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
		Properties: &workspaces.WorkspaceProperties{
			Sku:                             sku,
			PublicNetworkAccessForIngestion: pointer.ToEnum[workspaces.PublicNetworkAccessType](d.Get("internet_ingestion_access_type").(string)),
			PublicNetworkAccessForQuery:     pointer.ToEnum[workspaces.PublicNetworkAccessType](d.Get("internet_query_access_type").(string)),
			RetentionInDays:                 pointer.To(int64(d.Get("retention_in_days").(int))),
			Features: &workspaces.WorkspaceFeatures{
				EnableLogAccessUsingOnlyResourcePermissions: pointer.To(allowResourceOnlyPermission),
				DisableLocalAuth: pointer.To(!d.Get("local_authentication_enabled").(bool)),
			},
		},
	}

	// nolint : staticcheck
	if v, ok := d.GetOkExists("cmk_for_query_forced"); ok {
		parameters.Properties.ForceCmkForQuery = pointer.To(v.(bool))
	}

	if dailyQuotaGb, ok := d.GetOk("daily_quota_gb"); ok {
		parameters.Properties.WorkspaceCapping = &workspaces.WorkspaceCapping{
			DailyQuotaGb: pointer.To(dailyQuotaGb.(float64)),
		}
	}

	// The `ImmediatePurgeDataOn30Days` are not returned before it has been set
	// nolint : staticcheck
	if v, ok := d.GetOkExists("immediate_data_purge_on_30_days_enabled"); ok {
		parameters.Properties.Features.ImmediatePurgeDataOn30Days = pointer.To(v.(bool))
	}

	propName := "reservation_capacity_in_gb_per_day"
	capacityReservationLevel, ok := d.GetOk(propName)
	if ok {
		if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			parameters.Properties.Sku.CapacityReservationLevel = pointer.To(int64(capacityReservationLevel.(int)))
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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return err
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// `data_collection_rule_id` also needs an additional update.
	// error message: Default dcr is not applicable on workspace creation, please provide it on update.
	if v, ok := d.GetOk("data_collection_rule_id"); ok {
		parameters.Properties.DefaultDataCollectionRuleResourceId = pointer.To(v.(string))
	}

	// `allow_resource_only_permissions` needs an additional update, tracked in https://github.com/Azure/azure-rest-api-specs/issues/21591
	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
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
				return resp, "error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Features != nil && resp.Model.Properties.Features.EnableLogAccessUsingOnlyResourcePermissions != nil {
				return resp, strconv.FormatBool(*resp.Model.Properties.Features.EnableLogAccessUsingOnlyResourcePermissions), nil
			}

			return resp, "false", fmt.Errorf("retrieving %s: feature is nil", id)
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting on update for %s: %+v", id, err)
	}

	return resourceLogAnalyticsWorkspaceRead(d, meta)
}

func resourceLogAnalyticsWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspaceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}
	payload := existing.Model
	props := payload.Properties

	if props.Features == nil {
		props.Features = &workspaces.WorkspaceFeatures{}
	}

	if d.HasChange("allow_resource_only_permissions") {
		props.Features.EnableLogAccessUsingOnlyResourcePermissions = pointer.To(d.Get("allow_resource_only_permissions").(bool))
	}

	if d.HasChange("local_authentication_enabled") {
		props.Features.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if d.HasChange("cmk_for_query_forced") {
		props.ForceCmkForQuery = pointer.To(d.Get("cmk_for_query_forced").(bool))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`")
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("internet_ingestion_access_type") {
		props.PublicNetworkAccessForIngestion = pointer.ToEnum[workspaces.PublicNetworkAccessType](d.Get("internet_ingestion_access_type").(string))
	}

	if d.HasChange("internet_query_access_type") {
		props.PublicNetworkAccessForQuery = pointer.ToEnum[workspaces.PublicNetworkAccessType](d.Get("internet_query_access_type").(string))
	}

	var isLACluster bool
	if d.HasChange("sku") {
		skuName := d.Get("sku").(string)
		if sku := props.Sku; sku != nil {
			if strings.EqualFold(string(sku.Name), string(workspaces.WorkspaceSkuNameEnumLACluster)) {
				isLACluster = true
			}
		}

		if props.Sku == nil {
			props.Sku = &workspaces.WorkspaceSku{}
		}

		switch {
		case isLACluster:
			props.Sku.Name = workspaces.WorkspaceSkuNameEnumLACluster
		case skuName == "":
			// Default value if sku is not defined
			props.Sku.Name = workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight
		default:
			props.Sku.Name = workspaces.WorkspaceSkuNameEnum(skuName)
		}
	}

	if d.HasChange("reservation_capacity_in_gb_per_day") {
		skuName := d.Get("sku").(string)
		if payload.Properties.Sku == nil {
			payload.Properties.Sku = &workspaces.WorkspaceSku{}
		}

		if capacityReservationLevel, ok := d.GetOk("reservation_capacity_in_gb_per_day"); ok {
			if !strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
				return errors.New("`reservation_capacity_in_gb_per_day` can only be used with the `CapacityReservation` SKU")
			}

			props.Sku.CapacityReservationLevel = pointer.To(int64(capacityReservationLevel.(int)))
		} else if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			return errors.New("`reservation_capacity_in_gb_per_day` must be set when using the `CapacityReservation` SKU")
		}
	}

	if d.HasChange("retention_in_days") {
		props.RetentionInDays = pointer.To(int64(d.Get("retention_in_days").(int)))
	}

	if d.HasChange("daily_quota_gb") {
		props.WorkspaceCapping = &workspaces.WorkspaceCapping{
			DailyQuotaGb: pointer.To(d.Get("daily_quota_gb").(float64)),
		}
	}

	if d.HasChange("data_collection_rule_id") {
		props.DefaultDataCollectionRuleResourceId = pointer.To(d.Get("data_collection_rule_id").(string))
	}

	if d.HasChange("immediate_data_purge_on_30_days_enabled") {
		if props.Features == nil {
			props.Features = &workspaces.WorkspaceFeatures{}
		}
		props.Features.ImmediatePurgeDataOn30Days = pointer.To(d.Get("immediate_data_purge_on_30_days_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = expandTags(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

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
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceLogAnalyticsWorkspaceFlatten(ctx, sharedKeyClient, d, id, resp.Model, true)
}

func resourceLogAnalyticsWorkspaceFlatten(ctx context.Context, sharedKeyClient *sharedKeyWorkspaces.WorkspacesClient, d *pluginsdk.ResourceData, id *workspaces.WorkspaceId, model *workspaces.Workspace, includeResource bool) error {
	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		if model.Identity != nil {
			flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			d.Set("identity", flattenedIdentity)
		}

		if props := model.Properties; props != nil {
			d.Set("internet_ingestion_access_type", string(pointer.From(props.PublicNetworkAccessForIngestion)))
			d.Set("internet_query_access_type", string(pointer.From(props.PublicNetworkAccessForQuery)))

			d.Set("workspace_id", pointer.From(props.CustomerId))

			var skuName string
			if sku := props.Sku; sku != nil {
				skuName = string(sku.Name)
				d.Set("sku", skuName)

				if capacityReservationLevel := sku.CapacityReservationLevel; capacityReservationLevel != nil {
					d.Set("reservation_capacity_in_gb_per_day", pointer.From(capacityReservationLevel))
				}
			}

			d.Set("cmk_for_query_forced", pointer.From(props.ForceCmkForQuery))

			d.Set("retention_in_days", pointer.From(props.RetentionInDays))

			switch {
			case props.WorkspaceCapping != nil && props.WorkspaceCapping.DailyQuotaGb != nil:
				d.Set("daily_quota_gb", props.WorkspaceCapping.DailyQuotaGb)
			default:
				d.Set("daily_quota_gb", pointer.To(-1))
			}

			allowResourceOnlyPermissions := true
			purgeDataOnThirtyDays := false
			if lawFeatures := props.Features; lawFeatures != nil {
				allowResourceOnlyPermissions = pointer.From(lawFeatures.EnableLogAccessUsingOnlyResourcePermissions)
				if lawFeatures.DisableLocalAuth != nil {
					d.Set("local_authentication_enabled", !pointer.From(lawFeatures.DisableLocalAuth))
				}
				purgeDataOnThirtyDays = pointer.From(lawFeatures.ImmediatePurgeDataOn30Days)
			}
			d.Set("allow_resource_only_permissions", allowResourceOnlyPermissions)
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

			if includeResource {
				sharedKeyId := sharedKeyWorkspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
				sharedKeysResp, err := sharedKeyClient.SharedKeysGetSharedKeys(ctx, sharedKeyId)
				if err != nil {
					log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", id.WorkspaceName, err)
				} else {
					if sharedKeysModel := sharedKeysResp.Model; sharedKeysModel != nil {
						d.Set("primary_shared_key", pointer.From(sharedKeysModel.PrimarySharedKey))
						d.Set("secondary_shared_key", pointer.From(sharedKeysModel.SecondarySharedKey))
					}
				}
			}
		}

		d.Set("location", location.Normalize(model.Location))

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceLogAnalyticsWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SharedKeyWorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}
	sharedKeyId := sharedKeyWorkspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	permanentlyDeleteOnDestroy := meta.(*clients.Client).Features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy
	if err = client.DeleteThenPoll(ctx, sharedKeyId, sharedKeyWorkspaces.DeleteOperationOptions{Force: pointer.To(permanentlyDeleteOnDestroy)}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
