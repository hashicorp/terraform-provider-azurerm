// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenVpnGatewayConnectionVpnSiteLinkConnections(input *[]virtualwans.VpnSiteLinkConnection) any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, item := range *input {
		if item.Properties == nil {
			continue
		}

		props := *item.Properties

		connectionProtocolType := ""
		if props.VpnConnectionProtocolType != nil {
			connectionProtocolType = string(*props.VpnConnectionProtocolType)
		}

		vpnLinkConnectionMode := ""
		if props.VpnLinkConnectionMode != nil {
			vpnLinkConnectionMode = string(*props.VpnLinkConnectionMode)
		}

		vpnSiteLinkId := ""
		if props.VpnSiteLink != nil && props.VpnSiteLink.Id != nil {
			vpnSiteLinkId = *props.VpnSiteLink.Id
		}

		output = append(output, map[string]any{
			"name":                                  pointer.From(item.Name),
			"dpd_timeout_seconds":                   int(pointer.From(props.DpdTimeoutSeconds)),
			"egress_nat_rule_ids":                   flattenVpnGatewayConnectionNatRuleIds(props.EgressNatRules),
			"ingress_nat_rule_ids":                  flattenVpnGatewayConnectionNatRuleIds(props.IngressNatRules),
			"vpn_site_link_id":                      vpnSiteLinkId,
			"route_weight":                          int(pointer.From(props.RoutingWeight)),
			"protocol":                              connectionProtocolType,
			"connection_mode":                       vpnLinkConnectionMode,
			"bandwidth_mbps":                        int(pointer.From(props.ConnectionBandwidth)),
			"shared_key":                            pointer.From(props.SharedKey),
			"bgp_enabled":                           pointer.From(props.EnableBgp),
			"ipsec_policy":                          flattenVpnGatewayConnectionIpSecPolicies(props.IPsecPolicies),
			"ratelimit_enabled":                     pointer.From(props.EnableRateLimiting),
			"local_azure_ip_address_enabled":        pointer.From(props.UseLocalAzureIPAddress),
			"policy_based_traffic_selector_enabled": pointer.From(props.UsePolicyBasedTrafficSelectors),
			"custom_bgp_address":                    flattenVpnGatewayConnectionCustomBgpAddresses(props.VpnGatewayCustomBgpAddresses),
		})
	}

	return output
}

func flattenVpnGatewayConnectionIpSecPolicies(input *[]virtualwans.IPsecPolicy) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)
	for _, item := range *input {
		output = append(output, map[string]any{
			"sa_lifetime_sec":          int(item.SaLifeTimeSeconds),
			"sa_data_size_kb":          int(item.SaDataSizeKilobytes),
			"encryption_algorithm":     string(item.IPsecEncryption),
			"integrity_algorithm":      string(item.IPsecIntegrity),
			"ike_encryption_algorithm": string(item.IkeEncryption),
			"ike_integrity_algorithm":  string(item.IkeIntegrity),
			"dh_group":                 string(item.DhGroup),
			"pfs_group":                string(item.PfsGroup),
		})
	}

	return output
}

func flattenVpnGatewayConnectionRoutingConfiguration(input *virtualwans.RoutingConfiguration) []any {
	if input == nil {
		return []any{}
	}

	associateRouteTable := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associateRouteTable = *input.AssociatedRouteTable.Id
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
			"propagated_route_table": flattenVpnGatewayConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"associated_route_table": associateRouteTable,
			"inbound_route_map_id":   inboundRouteMapId,
			"outbound_route_map_id":  outboundRouteMapId,
		},
	}
}

func flattenVpnGatewayConnectionPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []any {
	if input == nil {
		return make([]any, 0)
	}

	labels := make([]any, 0)
	if input.Labels != nil {
		labels = pluginsdk.FlattenSlice(input.Labels)
	}

	routeTableIds := make([]any, 0)
	if input.Ids != nil {
		for _, id := range *input.Ids {
			if id.Id == nil {
				continue
			}
			routeTableIds = append(routeTableIds, *id.Id)
		}
	}

	return []any{
		map[string]any{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenVpnGatewayConnectionTrafficSelectorPolicy(input *[]virtualwans.TrafficSelectorPolicy) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"local_address_ranges":  pluginsdk.FlattenSlice(&item.LocalAddressRanges),
			"remote_address_ranges": pluginsdk.FlattenSlice(&item.RemoteAddressRanges),
		})
	}

	return results
}

func flattenVpnGatewayConnectionNatRuleIds(input *[]virtualwans.SubResource) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		id := pointer.From(item.Id)

		results = append(results, id)
	}

	return results
}

func flattenVpnGatewayConnectionCustomBgpAddresses(input *[]virtualwans.GatewayCustomBgpIPAddressIPConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"ip_address":          item.CustomBgpIPAddress,
			"ip_configuration_id": item.IPConfigurationId,
		})
	}

	return results
}
