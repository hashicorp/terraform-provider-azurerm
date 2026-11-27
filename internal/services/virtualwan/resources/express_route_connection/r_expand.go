// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandExpressRouteConnectionRouting(input []any) *expressrouteconnections.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return &expressrouteconnections.RoutingConfiguration{}
	}

	v := input[0].(map[string]any)
	result := expressrouteconnections.RoutingConfiguration{}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &expressrouteconnections.SubResource{
			Id: pointer.To(associatedRouteTableId),
		}
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		result.InboundRouteMap = &expressrouteconnections.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		result.OutboundRouteMap = &expressrouteconnections.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	if propagatedRouteTable := v["propagated_route_table"].([]any); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandExpressRouteConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return &result
}

func expandExpressRouteConnectionPropagatedRouteTable(input []any) *expressrouteconnections.PropagatedRouteTable {
	if len(input) == 0 || input[0] == nil {
		return &expressrouteconnections.PropagatedRouteTable{}
	}

	v := input[0].(map[string]any)

	result := expressrouteconnections.PropagatedRouteTable{}

	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = pluginsdk.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]any); len(routeTableIds) != 0 {
		result.Ids = expandExpressRouteIDsToSubResources(routeTableIds)
	}

	return &result
}

func expandExpressRouteIDsToSubResources(input []any) *[]expressrouteconnections.SubResource {
	ids := make([]expressrouteconnections.SubResource, 0)

	for _, v := range input {
		ids = append(ids, expressrouteconnections.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}
