// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutegateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRouteConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ExpressRouteConnectionName,
		},

		"express_route_circuit_peering_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
		},

		"express_route_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: expressroutegateways.ValidateExpressRouteGatewayID,
		},

		"authorization_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"internet_security_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"express_route_gateway_bypass_enabled": {
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

					"associated_route_table_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: virtualwans.ValidateHubRouteTableID,
						AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
					},

					"propagated_route_table": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"labels": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
								},

								"route_table_ids": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: virtualwans.ValidateHubRouteTableID,
									},
									AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
								},
							},
						},
						AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
					},
				},
			},
		},

		"routing_weight": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 32000),
		},
	}
}
