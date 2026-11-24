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
			Description:  "The type of the ExpressRoute Circuit Peering. Acceptable values include `AzurePrivatePeering`, `AzurePublicPeering` and `MicrosoftPeering`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(expressroutecircuitpeerings.PossibleValuesForExpressRoutePeeringType(), false),
		},

		"express_route_circuit_name": {
			Description: "The name of the ExpressRoute Circuit in which to create the Peering.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"primary_peer_address_prefix": {
			Description: "A `/30` subnet for the primary link.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			RequiredWith: []string{
				"secondary_peer_address_prefix",
			},
		},

		"secondary_peer_address_prefix": {
			Description: "A `/30` subnet for the secondary link.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			RequiredWith: []string{
				"primary_peer_address_prefix",
			},
		},

		"ipv4_enabled": {
			Description: "A boolean value indicating whether the IPv4 peering is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"vlan_id": {
			Description: "A valid VLAN ID to establish this peering on.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"shared_key": {
			Description:  "The shared key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(1, 25),
		},

		"peer_asn": {
			Description: "Either a 16-bit or a 32-bit ASN. Can either be public or private.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
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
			Description: "The ASN used by Azure.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"primary_azure_port": {
			Description: "The Primary Port used by Azure for this Peering.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"secondary_azure_port": {
			Description: "The Secondary Port used by Azure for this Peering.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"route_filter_id": {
			Description:  "The ID of the Route Filter. Only available when `peering_type` is set to `MicrosoftPeering`.",
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
