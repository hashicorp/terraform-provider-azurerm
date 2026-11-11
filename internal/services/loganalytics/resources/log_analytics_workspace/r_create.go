// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsWorkspaceCreate(d *pluginsdk.ResourceData, meta any) error {
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
		Tags:     expandTags(d.Get("tags").(map[string]any)),
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

	//nolint:staticcheck
	if v, ok := d.GetOkExists("cmk_for_query_forced"); ok {
		parameters.Properties.ForceCmkForQuery = pointer.To(v.(bool))
	}

	if dailyQuotaGb, ok := d.GetOk("daily_quota_gb"); ok {
		parameters.Properties.WorkspaceCapping = &workspaces.WorkspaceCapping{
			DailyQuotaGb: pointer.To(dailyQuotaGb.(float64)),
		}
	}

	// The `ImmediatePurgeDataOn30Days` are not returned before it has been set
	//nolint:staticcheck
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
		expanded, err := identity.ExpandSystemOrUserAssignedMap(v.([]any))
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
		Refresh: func() (any, string, error) {
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
