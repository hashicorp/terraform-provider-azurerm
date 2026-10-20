// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualNetworkGatewayConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute),
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec),
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet),
			}, false),
		},

		"virtual_network_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},

		"shared_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C the API generates a key for the user if not supplied
			Computed:  true,
			Sensitive: true,
		},

		"authorization_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dpd_timeout_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
		},

		"express_route_circuit_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: expressroutecircuits.ValidateExpressRouteCircuitID,
		},

		"egress_nat_rule_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgateways.ValidateVirtualNetworkGatewayNatRuleID,
			},
		},

		"ingress_nat_rule_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgateways.ValidateVirtualNetworkGatewayNatRuleID,
			},
		},

		"peer_virtual_network_gateway_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},

		"local_azure_ip_address_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"local_network_gateway_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: localnetworkgateways.ValidateLocalNetworkGatewayID,
		},

		"bgp_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"use_policy_based_traffic_selectors": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"routing_weight": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 32000),
		},

		"express_route_gateway_bypass": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"private_link_fast_path_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"connection_protocol": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForVirtualNetworkGatewayConnectionProtocol(), false),
		},

		"connection_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForVirtualNetworkGatewayConnectionMode(), false),
			Default:      string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeDefault),
		},

		"traffic_selector_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"local_address_cidrs": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"remote_address_cidrs": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"custom_bgp_addresses": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"primary": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.Any(validation.IsIPv4Address),
					},
					"secondary": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.Any(validation.IsIPv4Address),
					},
				},
			},
		},

		"ipsec_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dh_group": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForDhGroup(), false),
					},

					"ike_encryption": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIkeEncryption(), false),
					},

					"ike_integrity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIkeIntegrity(), false),
					},

					"ipsec_encryption": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIPsecEncryption(), false),
					},

					"ipsec_integrity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIPsecIntegrity(), false),
					},

					"pfs_group": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForPfsGroup(), false),
					},

					"sa_datasize": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},

					"sa_lifetime": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntAtLeast(300),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
