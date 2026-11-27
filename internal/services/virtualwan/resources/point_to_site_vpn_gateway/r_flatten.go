// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenPointToSiteVPNGatewayConnectionConfiguration(input *[]virtualwans.P2SConnectionConfiguration) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		route := make([]any, 0)
		addressPrefixes := make([]any, 0)
		enableInternetSecurity := false
		if props := v.Properties; props != nil {
			if props.VpnClientAddressPool == nil {
				continue
			}

			if props.VpnClientAddressPool.AddressPrefixes != nil {
				for _, prefix := range *props.VpnClientAddressPool.AddressPrefixes {
					addressPrefixes = append(addressPrefixes, prefix)
				}
			}

			if props.EnableInternetSecurity != nil {
				enableInternetSecurity = *props.EnableInternetSecurity
			}

			if props.RoutingConfiguration != nil {
				route = flattenPointToSiteVPNGatewayConnectionRouteConfiguration(props.RoutingConfiguration)
			}
		}

		output = append(output, map[string]any{
			"name": pointer.From(v.Name),
			"vpn_client_address_pool": []any{
				map[string]any{
					"address_prefixes": addressPrefixes,
				},
			},
			"route":                     route,
			"internet_security_enabled": enableInternetSecurity,
		})
	}

	return output
}

func flattenPointToSiteVPNGatewayConnectionRouteConfiguration(input *virtualwans.RoutingConfiguration) []any {
	if input == nil {
		return []any{}
	}

	var associatedRouteTableId string
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
	}

	var inboundRouteMapId string
	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		inboundRouteMapId = *input.InboundRouteMap.Id
	}

	var outboundRouteMapId string
	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		outboundRouteMapId = *input.OutboundRouteMap.Id
	}

	return []any{
		map[string]any{
			"associated_route_table_id": associatedRouteTableId,
			"inbound_route_map_id":      inboundRouteMapId,
			"outbound_route_map_id":     outboundRouteMapId,
			"propagated_route_table":    flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input.PropagatedRouteTables),
		},
	}
}

func flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []any {
	if input == nil {
		return []any{}
	}
	ids := make([]string, 0)
	if input.Ids != nil {
		for _, item := range *input.Ids {
			if item.Id != nil {
				ids = append(ids, *item.Id)
			}
		}
	}
	return []any{
		map[string]any{
			"ids":    ids,
			"labels": pluginsdk.FlattenSlice(input.Labels),
		},
	}
}
