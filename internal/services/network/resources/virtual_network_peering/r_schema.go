// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualNetworkPeeringSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Virtual Network peering.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"virtual_network_name": {
			Description: "The name of the Virtual Network.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"remote_virtual_network_id": {
			Description:  "The full Azure resource ID of the remote Virtual Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"allow_virtual_network_access": {
			Description: "Controls if the traffic from the local Virtual Network can reach the remote Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"allow_forwarded_traffic": {
			Description: "Controls if forwarded traffic from VMs in the remote Virtual Network is allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"allow_gateway_transit": {
			Description: "Controls gatewayLinks can be used in the remote Virtual Network's link to the local Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"local_subnet_names": {
			Description: "A list of local Subnet names that are Subnet peered with remote Virtual Network.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"only_ipv6_peering_enabled": {
			Description: "Whether only IPv6 address space is peered for Subnet peering.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"peer_complete_virtual_networks_enabled": {
			Description: "Whether complete Virtual Network address space is peered.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},

		"remote_subnet_names": {
			Description: "A list of remote Subnet names from remote Virtual Network that are Subnet peered.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"use_remote_gateways": {
			Description: "Controls if remote gateways can be used on the local Virtual Network. If the flag is set to `true`, and `allow_gateway_transit` on the remote peering is also `true`, Virtual Network will use gateways of remote Virtual Network for transit. Only one peering can have this flag set to `true`. This flag cannot be set if Virtual Network already has a gateway.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"triggers": {
			Description: "A mapping of key values pairs that can be used to sync network routes from the remote Virtual Network to the local Virtual Network. See [the trigger example](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network_peering#example-usage-triggers) for an example on how to set it up.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
