// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualHubConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualHubConnectionName,
		},

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVirtualHubID,
		},

		"remote_virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"internet_security_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"associated_route_table_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: virtualwans.ValidateHubRouteTableID,
						AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
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
						Type:         pluginsdk.TypeList,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						MaxItems:     1,
						AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"labels": {
									Type:         pluginsdk.TypeSet,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"route_table_ids": {
									Type:         pluginsdk.TypeList,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: virtualwans.ValidateHubRouteTableID,
									},
								},
							},
						},
					},

					"static_vnet_local_route_override_criteria": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(virtualwans.VnetLocalRouteOverrideCriteriaContains),
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVnetLocalRouteOverrideCriteria(), false),
					},

					"static_vnet_propagate_static_routes_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					// lintignore:XS003
					"static_vnet_route": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"address_prefixes": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsCIDR,
									},
								},

								"next_hop_ip_address": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},
				},
			},
		},
	}
}
