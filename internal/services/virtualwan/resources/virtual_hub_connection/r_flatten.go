// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceVirtualHubConnectionFlatten(d *pluginsdk.ResourceData, id *virtualwans.HubVirtualNetworkConnectionId, model *virtualwans.HubVirtualNetworkConnection) error {
	d.Set("name", id.HubVirtualNetworkConnectionName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("internet_security_enabled", props.EnableInternetSecurity)
			remoteVirtualNetworkId := ""
			if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.Id != nil {
				parsedRemoteVirtualNetworkId, err := commonids.ParseVirtualNetworkIDInsensitively(*props.RemoteVirtualNetwork.Id)
				if err != nil {
					return err
				}
				remoteVirtualNetworkId = parsedRemoteVirtualNetworkId.ID()
			}
			d.Set("remote_virtual_network_id", remoteVirtualNetworkId)

			if err := d.Set("routing", flattenVirtualHubConnectionRouting(props.RoutingConfiguration)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenVirtualHubConnectionRouting(input *virtualwans.RoutingConfiguration) []any {
	if input == nil {
		return []any{}
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
	}

	inboundRouteMapId := ""
	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		inboundRouteMapId = *input.InboundRouteMap.Id
	}

	outboundRouteMapId := ""
	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		outboundRouteMapId = *input.OutboundRouteMap.Id
	}

	staticVnetLocalRouteOverrideCriteria := ""
	if input.VnetRoutes != nil && input.VnetRoutes.StaticRoutesConfig != nil && input.VnetRoutes.StaticRoutesConfig.VnetLocalRouteOverrideCriteria != nil {
		staticVnetLocalRouteOverrideCriteria = string(*input.VnetRoutes.StaticRoutesConfig.VnetLocalRouteOverrideCriteria)
	}

	staticVnetPropagateStaticRoutes := true
	if input.VnetRoutes != nil && input.VnetRoutes.StaticRoutesConfig != nil {
		staticVnetPropagateStaticRoutes = pointer.From(input.VnetRoutes.StaticRoutesConfig.PropagateStaticRoutes)
	}

	return []any{
		map[string]any{
			"associated_route_table_id":                   associatedRouteTableId,
			"inbound_route_map_id":                        inboundRouteMapId,
			"outbound_route_map_id":                       outboundRouteMapId,
			"propagated_route_table":                      flattenVirtualHubConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"static_vnet_route":                           flattenVirtualHubConnectionVnetStaticRoute(input.VnetRoutes),
			"static_vnet_local_route_override_criteria":   staticVnetLocalRouteOverrideCriteria,
			"static_vnet_propagate_static_routes_enabled": staticVnetPropagateStaticRoutes,
		},
	}
}

func flattenVirtualHubConnectionPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []any {
	if input == nil {
		return make([]any, 0)
	}

	labels := make([]any, 0)
	if input.Labels != nil {
		labels = pluginsdk.FlattenSlice(input.Labels)
	}

	routeTableIds := make([]any, 0)
	if input.Ids != nil {
		routeTableIds = flattenVirtualWANSubResourcesToIDs(input.Ids)
	}

	return []any{
		map[string]any{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenVirtualHubConnectionVnetStaticRoute(input *virtualwans.VnetRoute) []any {
	results := make([]any, 0)
	if input == nil || input.StaticRoutes == nil {
		return results
	}

	for _, item := range *input.StaticRoutes {
		addressPrefixes := make([]any, 0)
		if item.AddressPrefixes != nil {
			addressPrefixes = pluginsdk.FlattenSlice(item.AddressPrefixes)
		}

		v := map[string]any{
			"name":                pointer.From(item.Name),
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": pointer.From(item.NextHopIPAddress),
		}

		results = append(results, v)
	}

	return results
}

func flattenVirtualWANSubResourcesToIDs(input *[]virtualwans.SubResource) []any {
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
