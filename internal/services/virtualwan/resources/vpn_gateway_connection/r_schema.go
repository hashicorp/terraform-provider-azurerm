// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"math"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vpnGatewayConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vpn_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVpnGatewayID,
		},

		"remote_vpn_site_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVpnSiteID,
		},

		"internet_security_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// Note: O+C because the service will create a route table for the user if this is not specified.
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"associated_route_table": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: virtualwans.ValidateHubRouteTableID,
					},

					"inbound_route_map_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: virtualwans.ValidateRouteMapID,
					},

					"outbound_route_map_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: virtualwans.ValidateRouteMapID,
					},

					"propagated_route_table": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"route_table_ids": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: virtualwans.ValidateHubRouteTableID,
									},
								},

								"labels": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"vpn_link": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"vpn_site_link_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// The vpn site link associated with one link connection can not be updated
						ForceNew:     true,
						ValidateFunc: virtualwans.ValidateVpnSiteLinkID,
					},

					"dpd_timeout_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(9, 3600),
					},

					"egress_nat_rule_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: virtualwans.ValidateNatRuleID,
						},
					},

					"ingress_nat_rule_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: virtualwans.ValidateNatRuleID,
						},
					},

					"connection_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnLinkConnectionMode(), false),
						Default:      string(virtualwans.VpnLinkConnectionModeDefault),
					},

					"route_weight": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						Default:      0,
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVirtualNetworkGatewayConnectionProtocol(), false),
						Default:      string(virtualwans.VirtualNetworkGatewayConnectionProtocolIKEvTwo),
					},

					"bandwidth_mbps": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      10,
					},

					"shared_key": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C the API generates a key for the user if not supplied
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"bgp_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  false,
					},

					"ipsec_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"sa_lifetime_sec": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(300, 172799),
								},
								"sa_data_size_kb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, math.MaxInt32),
								},
								"encryption_algorithm": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecEncryption(), false),
								},
								"integrity_algorithm": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecIntegrity(), false),
								},

								"ike_encryption_algorithm": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeEncryption(), false),
								},

								"ike_integrity_algorithm": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeIntegrity(), false),
								},

								"dh_group": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForDhGroup(), false),
								},

								"pfs_group": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForPfsGroup(), false),
								},
							},
						},
					},

					"ratelimit_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"local_azure_ip_address_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"policy_based_traffic_selector_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"custom_bgp_address": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"ip_configuration_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"traffic_selector_policy": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"local_address_ranges": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDR,
						},
					},

					"remote_address_ranges": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDR,
						},
					},
				},
			},
		},
	}
}
