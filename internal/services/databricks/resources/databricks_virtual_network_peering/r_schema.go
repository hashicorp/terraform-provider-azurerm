// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func databricksVirtualNetworkPeeringSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabricksVirtualNetworkPeeringName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"remote_address_space_prefixes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,

			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDRIsIPv4OrIPv6,
			},
		},

		"remote_virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"address_space_prefixes": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"allow_virtual_network_access": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"allow_forwarded_traffic": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"allow_gateway_transit": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"use_remote_gateways": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}
