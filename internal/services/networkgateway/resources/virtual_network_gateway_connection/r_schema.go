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
			Description:  "The name of the connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"type": {
			Description: "The type of connection. Each connection type requires different mandatory arguments (refer to the examples above).",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute),
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec),
				string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet),
			}, false),
		},

		"virtual_network_gateway_id": {
			Description:  "The ID of the Virtual Network Gateway in which the connection will be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},

		"shared_key": {
			Description: "The shared IPSec key. A key could be provided if a Site-to-Site, VNet-to-VNet or ExpressRoute connection is created.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C the API generates a key for the user if not supplied
			Computed:  true,
			Sensitive: true,
		},

		"authorization_key": {
			Description:  "The authorization key associated with the Express Route Circuit. This field is required only if the type is an ExpressRoute connection.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dpd_timeout_seconds": {
			Description: "The dead peer detection timeout of this connection in seconds.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			ForceNew:    true,
		},

		"express_route_circuit_id": {
			Description:  "The ID of the Express Route Circuit when creating an ExpressRoute connection (i.e. when `type` is `ExpressRoute`). The Express Route Circuit can be in the same or in a different subscription.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: expressroutecircuits.ValidateExpressRouteCircuitID,
		},

		"egress_nat_rule_ids": {
			Description: "A list of the egress NAT Rule IDs.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgateways.ValidateVirtualNetworkGatewayNatRuleID,
			},
		},

		"ingress_nat_rule_ids": {
			Description: "A list of the ingress NAT Rule IDs.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgateways.ValidateVirtualNetworkGatewayNatRuleID,
			},
		},

		"peer_virtual_network_gateway_id": {
			Description:  "The ID of the peer Virtual Network Gateway when creating a VNet-to-VNet connection (i.e. when `type` is `Vnet2Vnet`). The peer Virtual Network Gateway can be in the same or in a different subscription.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},

		"local_azure_ip_address_enabled": {
			Description: "Use private local Azure IP for the connection.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"local_network_gateway_id": {
			Description:  "The ID of the Local Network Gateway when creating Site-to-Site connection (i.e. when `type` is `IPsec`).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: localnetworkgateways.ValidateLocalNetworkGatewayID,
		},

		"bgp_enabled": {
			Description: "If `true`, BGP (Border Gateway Protocol) is enabled for this connection.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"use_policy_based_traffic_selectors": {
			Description: "If `true`, policy-based traffic selectors are enabled for this connection. Enabling policy-based traffic selectors requires an `ipsec_policy` block.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"routing_weight": {
			Description:  "The routing weight.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 32000),
		},

		"express_route_gateway_bypass": {
			Description: "If `true`, data packets will bypass ExpressRoute Gateway for data forwarding. This is only valid for ExpressRoute connections.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"private_link_fast_path_enabled": {
			Description: "Bypass the Express Route gateway when accessing private-links. When enabled `express_route_gateway_bypass` must be set to `true`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"connection_protocol": {
			Description:  "The IKE protocol version to use.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForVirtualNetworkGatewayConnectionProtocol(), false),
		},

		"connection_mode": {
			Description:  "Connection mode to use.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForVirtualNetworkGatewayConnectionMode(), false),
			Default:      string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeDefault),
		},

		"traffic_selector_policy": {
			Description: "A `traffic_selector_policy` allows to specify a traffic selector policy proposal to be used in a Virtual Network Gateway connection. For details about traffic selectors refer to [the relevant section in the Azure documentation](https://learn.microsoft.com/azure/vpn-gateway/vpn-gateway-connect-multiple-policybased-rm-ps).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"local_address_cidrs": {
						Description: "List of local CIDRs.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"remote_address_cidrs": {
						Description: "List of remote CIDRs.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"custom_bgp_addresses": {
			Description: "The block can only be used on `IPSec` / `active-active` connections. For details see [the relevant section in the Azure documentation](https://learn.microsoft.com/azure/vpn-gateway/vpn-gateway-howto-aws-bgp).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"primary": {
						Description:  "Single IP address that is part of the `azurerm_virtual_network_gateway` ip_configuration (first one).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.Any(validation.IsIPv4Address),
					},
					"secondary": {
						Description:  "Single IP address that is part of the `azurerm_virtual_network_gateway` ip_configuration (second one).",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.Any(validation.IsIPv4Address),
					},
				},
			},
		},

		"ipsec_policy": {
			Description: "Only a single policy can be defined for a connection. For details on custom policies refer to [the relevant section in the Azure documentation](https://learn.microsoft.com/azure/vpn-gateway/vpn-gateway-ipsecikepolicy-rm-powershell).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dh_group": {
						Description:  "The DH group used in IKE phase 1 for initial SA.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForDhGroup(), false),
					},

					"ike_encryption": {
						Description:  "The IKE encryption algorithm.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIkeEncryption(), false),
					},

					"ike_integrity": {
						Description:  "The IKE integrity algorithm.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIkeIntegrity(), false),
					},

					"ipsec_encryption": {
						Description:  "The IPSec encryption algorithm.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIPsecEncryption(), false),
					},

					"ipsec_integrity": {
						Description:  "The IPSec integrity algorithm.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForIPsecIntegrity(), false),
					},

					"pfs_group": {
						Description:  "The DH group used in IKE phase 2 for new child SA.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgatewayconnections.PossibleValuesForPfsGroup(), false),
					},

					"sa_datasize": {
						Description:  "The IPSec SA payload size in KB.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},

					"sa_lifetime": {
						Description:  "The IPSec SA lifetime in seconds.",
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
