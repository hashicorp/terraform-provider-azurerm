// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsWorkspaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
				50,
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
	}
}
