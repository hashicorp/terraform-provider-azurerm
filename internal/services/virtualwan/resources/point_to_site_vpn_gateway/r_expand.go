// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandPointToSiteVPNGatewayConnectionConfiguration(input []any) *[]virtualwans.P2SConnectionConfiguration {
	configurations := make([]virtualwans.P2SConnectionConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		addressPrefixes := make([]string, 0)
		name := raw["name"].(string)

		clientAddressPoolsRaw := raw["vpn_client_address_pool"].([]any)
		for _, clientV := range clientAddressPoolsRaw {
			clientRaw := clientV.(map[string]any)

			addressPrefixesRaw := clientRaw["address_prefixes"].(*pluginsdk.Set).List()
			for _, prefix := range addressPrefixesRaw {
				addressPrefixes = append(addressPrefixes, prefix.(string))
			}
		}

		configurations = append(configurations, virtualwans.P2SConnectionConfiguration{
			Name: pointer.To(name),
			Properties: &virtualwans.P2SConnectionConfigurationProperties{
				VpnClientAddressPool: &virtualwans.AddressSpace{
					AddressPrefixes: &addressPrefixes,
				},
				RoutingConfiguration:   expandPointToSiteVPNGatewayConnectionRouteConfiguration(raw["route"].([]any)),
				EnableInternetSecurity: pointer.To(raw["internet_security_enabled"].(bool)),
			},
		})
	}

	return &configurations
}

func expandPointToSiteVPNGatewayConnectionRouteConfiguration(input []any) *virtualwans.RoutingConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	routingConfiguration := &virtualwans.RoutingConfiguration{
		AssociatedRouteTable: &virtualwans.SubResource{
			Id: pointer.To(v["associated_route_table_id"].(string)),
		},
		PropagatedRouteTables: expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(v["propagated_route_table"].([]any)),
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		routingConfiguration.InboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		routingConfiguration.OutboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	return routingConfiguration
}

func expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input []any) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	idRaws := pluginsdk.ExpandStringSlice(v["ids"].([]any))
	ids := make([]virtualwans.SubResource, len(*idRaws))
	for i, item := range *idRaws {
		ids[i] = virtualwans.SubResource{
			Id: pointer.To(item),
		}
	}
	return &virtualwans.PropagatedRouteTable{
		Labels: pluginsdk.ExpandStringSlice(v["labels"].(*pluginsdk.Set).List()),
		Ids:    &ids,
	}
}
