// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenExpressRouteConnectionRouting(input *expressrouteconnections.RoutingConfiguration) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
	}
	routeTableId, err := virtualwans.ParseHubRouteTableIDInsensitively(associatedRouteTableId)
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"associated_route_table_id": routeTableId.ID(),
		"propagated_route_table":    flattenExpressRouteConnectionPropagatedRouteTable(input.PropagatedRouteTables),
	}

	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		result["inbound_route_map_id"] = input.InboundRouteMap.Id
	}

	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		result["outbound_route_map_id"] = input.OutboundRouteMap.Id
	}

	return []any{result}, nil
}

func flattenExpressRouteConnectionPropagatedRouteTable(input *expressrouteconnections.PropagatedRouteTable) []any {
	if input == nil {
		return make([]any, 0)
	}

	labels := make([]any, 0)
	if input.Labels != nil {
		labels = pluginsdk.FlattenSlice(input.Labels)
	}

	routeTableIds := make([]any, 0)
	if input.Ids != nil {
		routeTableIds = flattenExpressRouteSubResourcesToIDs(input.Ids)
	}

	return []any{
		map[string]any{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenExpressRouteSubResourcesToIDs(input *[]expressrouteconnections.SubResource) []any {
	ids := make([]any, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}
