// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vpnGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVirtualHubID,
		},

		"routing_preference": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Microsoft Network",
			ValidateFunc: validation.StringInSlice([]string{
				"Microsoft Network",
				"Internet",
			}, false),
		},

		"bgp_route_translation_for_nat_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"bgp_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"asn": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},

					"peer_weight": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},

					"bgp_peering_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"instance_0_bgp_peering_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_ips": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsIPv4Address,
									},
								},

								"ip_configuration_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"default_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"tunnel_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"instance_1_bgp_peering_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_ips": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsIPv4Address,
									},
								},

								"ip_configuration_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"default_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"tunnel_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"scale_unit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"tags": commonschema.Tags(),

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"public_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}
