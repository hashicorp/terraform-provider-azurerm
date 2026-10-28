// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsWorkspaceUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
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
		payload.Tags = expandTags(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceLogAnalyticsWorkspaceRead(d, meta)
}
