// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_peering

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuitpeerings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRouteCircuitPeeringSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"peering_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(expressroutecircuitpeerings.PossibleValuesForExpressRoutePeeringType(), false),
		},

		"express_route_circuit_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"primary_peer_address_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			RequiredWith: []string{
				"secondary_peer_address_prefix",
			},
		},

		"secondary_peer_address_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			RequiredWith: []string{
				"primary_peer_address_prefix",
			},
		},

		"ipv4_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"vlan_id": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"shared_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(1, 25),
		},

		"peer_asn": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"microsoft_peering_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"advertised_public_prefixes": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"customer_asn": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},

					"routing_registry_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "NONE",
					},

					"advertised_communities": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"ipv6": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"microsoft_peering": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"advertised_public_prefixes": {
									Type:     pluginsdk.TypeList,
									MinItems: 1,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsCIDR,
									},
								},

								"customer_asn": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  0,
								},

								"routing_registry_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "NONE",
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"advertised_communities": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"primary_peer_address_prefix": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"secondary_peer_address_prefix": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"route_filter_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: routefilters.ValidateRouteFilterID,
					},
				},
			},
		},

		"azure_asn": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"primary_azure_port": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_azure_port": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"route_filter_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: routefilters.ValidateRouteFilterID,
		},

		"gateway_manager_etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
