// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Databricks Virtual Network Peering resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabricksVirtualNetworkPeeringName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_id": {
			Description:  "The ID of the Databricks Workspace that this Databricks Virtual Network Peering is bound to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"remote_address_space_prefixes": {
			Description: "A list of address blocks reserved for the remote Virtual Network in CIDR notation.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,

			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDRIsIPv4OrIPv6,
			},
		},

		"remote_virtual_network_id": {
			Description:  "The ID of the remote Virtual Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"address_space_prefixes": {
			Description: "A list of address blocks reserved for this Virtual Network in CIDR notation.",
			Type:        pluginsdk.TypeList,
			Computed:    true,

			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_network_id": {
			Description: "The ID of the internal Virtual Network used by the Databricks Workspace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"allow_virtual_network_access": {
			Description: "Whether the VMs in the local Virtual Network space can access the VMs in the remote Virtual Network space.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"allow_forwarded_traffic": {
			Description: "Whether the forwarded traffic from the VMs in the local Virtual Network can be forwarded to the remote Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"allow_gateway_transit": {
			Description: "Whether the gateway links can be used in the remote Virtual Network to link to the Databricks Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"use_remote_gateways": {
			Description: "Whether remote gateways can be used on the Databricks Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}
