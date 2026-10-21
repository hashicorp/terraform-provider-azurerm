// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func networkWatcherFlowLogSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_watcher_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// lintignore: S013
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NetworkWatcherFlowLogName,
		},

		"target_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.Any(
				networksecuritygroups.ValidateNetworkSecurityGroupID,
				commonids.ValidateVirtualNetworkID,
				commonids.ValidateSubnetID,
				commonids.ValidateNetworkInterfaceID,
			),
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"retention_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:             pluginsdk.TypeBool,
						Required:         true,
						DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyEnabledDiff,
					},

					"days": {
						Type:             pluginsdk.TypeInt,
						Required:         true,
						DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyDaysDiff,
					},
				},
			},
		},

		"traffic_analytics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"workspace_region": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						StateFunc:        location.StateFunc,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"workspace_resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// TODO: check to see if empty values should be allowed (the previous validator permitted them)
						ValidateFunc: validation.Any(validation.StringIsEmpty, workspaces.ValidateWorkspaceID),
					},

					"interval_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntInSlice([]int{10, 60}),
						Default:      60,
					},
				},
			},
		},

		"version": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
			// NOTE: the Azure Flow Log API currently supports versions 1, 2 and 5 - this list should be
			// updated as the service adds support for additional versions.
			ValidateFunc: validation.IntInSlice([]int{1, 2, 5}),
		},

		"location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ForceNew:         true,
			ValidateFunc:     location.EnhancedValidate,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"tags": commonschema.Tags(),
	}
}
