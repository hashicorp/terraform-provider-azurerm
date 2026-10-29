// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	sharedKeyWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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
			d.Set("internet_ingestion_access_type", pointer.FromEnum(props.PublicNetworkAccessForIngestion))
			d.Set("internet_query_access_type", pointer.FromEnum(props.PublicNetworkAccessForQuery))

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
