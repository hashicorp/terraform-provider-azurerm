// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualHubConnectionRouting(input []any) *virtualwans.RoutingConfiguration {
	if len(input) == 0 {
		return &virtualwans.RoutingConfiguration{}
	}

	v := input[0].(map[string]any)

	result := &virtualwans.RoutingConfiguration{
		VnetRoutes: &virtualwans.VnetRoute{
			StaticRoutes: expandVirtualHubConnectionVnetStaticRoute(v["static_vnet_route"].([]any)),
			StaticRoutesConfig: &virtualwans.StaticRoutesConfig{
				PropagateStaticRoutes:          pointer.To(v["static_vnet_propagate_static_routes_enabled"].(bool)),
				VnetLocalRouteOverrideCriteria: pointer.ToEnum[virtualwans.VnetLocalRouteOverrideCriteria](v["static_vnet_local_route_override_criteria"].(string)),
			},
		},
	}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &virtualwans.SubResource{
			Id: pointer.To(associatedRouteTableId),
		}
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		result.InboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		result.OutboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	if propagatedRouteTable := v["propagated_route_table"].([]any); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandVirtualHubConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return result
}

func expandVirtualHubConnectionPropagatedRouteTable(input []any) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return &virtualwans.PropagatedRouteTable{}
	}

	v := input[0].(map[string]any)

	result := virtualwans.PropagatedRouteTable{}

	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = pluginsdk.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]any); len(routeTableIds) != 0 {
		result.Ids = expandIDsToVirtualWANSubResources(routeTableIds)
	}

	return &result
}

func expandVirtualHubConnectionVnetStaticRoute(input []any) *[]virtualwans.StaticRoute {
	if len(input) == 0 {
		return nil
	}

	results := make([]virtualwans.StaticRoute, 0)

	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)

		result := virtualwans.StaticRoute{}

		if name := v["name"].(string); name != "" {
			result.Name = pointer.To(name)
		}

		if addressPrefixes := v["address_prefixes"].(*pluginsdk.Set).List(); len(addressPrefixes) != 0 {
			result.AddressPrefixes = pluginsdk.ExpandStringSlice(addressPrefixes)
		}

		if nextHopIPAddress := v["next_hop_ip_address"].(string); nextHopIPAddress != "" {
			result.NextHopIPAddress = pointer.To(nextHopIPAddress)
		}

		results = append(results, result)
	}

	return &results
}

func expandIDsToVirtualWANSubResources(input []any) *[]virtualwans.SubResource {
	ids := make([]virtualwans.SubResource, 0)

	for _, v := range input {
		ids = append(ids, virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}
