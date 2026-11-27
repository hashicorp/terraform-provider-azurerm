// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVpnGatewayConnectionVpnSiteLinkConnections(input []any) *[]virtualwans.VpnSiteLinkConnection {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.VpnSiteLinkConnection, 0)
	for _, itemRaw := range input {
		item := itemRaw.(map[string]any)
		v := virtualwans.VpnSiteLinkConnection{
			Name: pointer.To(item["name"].(string)),
			Properties: &virtualwans.VpnSiteLinkConnectionProperties{
				VpnSiteLink: &virtualwans.SubResource{
					Id: pointer.To(item["vpn_site_link_id"].(string)),
				},
				RoutingWeight:                  pointer.To(int64(item["route_weight"].(int))),
				VpnConnectionProtocolType:      pointer.ToEnum[virtualwans.VirtualNetworkGatewayConnectionProtocol](item["protocol"].(string)),
				VpnLinkConnectionMode:          pointer.ToEnum[virtualwans.VpnLinkConnectionMode](item["connection_mode"].(string)),
				ConnectionBandwidth:            pointer.To(int64(item["bandwidth_mbps"].(int))),
				EnableBgp:                      pointer.To(item["bgp_enabled"].(bool)),
				IPsecPolicies:                  expandVpnGatewayConnectionIpSecPolicies(item["ipsec_policy"].([]any)),
				EnableRateLimiting:             pointer.To(item["ratelimit_enabled"].(bool)),
				UseLocalAzureIPAddress:         pointer.To(item["local_azure_ip_address_enabled"].(bool)),
				UsePolicyBasedTrafficSelectors: pointer.To(item["policy_based_traffic_selector_enabled"].(bool)),
				VpnGatewayCustomBgpAddresses:   expandVpnGatewayConnectionCustomBgpAddresses(item["custom_bgp_address"].(*pluginsdk.Set).List()),
			},
		}

		if dpdTimeoutSeconds := item["dpd_timeout_seconds"].(int); dpdTimeoutSeconds != 0 {
			v.Properties.DpdTimeoutSeconds = pointer.To(int64(dpdTimeoutSeconds))
		}

		if egressNatRuleIds := item["egress_nat_rule_ids"].(*pluginsdk.Set).List(); len(egressNatRuleIds) != 0 {
			v.Properties.EgressNatRules = expandVpnGatewayConnectionNatRuleIds(egressNatRuleIds)
		}

		if ingressNatRuleIds := item["ingress_nat_rule_ids"].(*pluginsdk.Set).List(); len(ingressNatRuleIds) != 0 {
			v.Properties.IngressNatRules = expandVpnGatewayConnectionNatRuleIds(ingressNatRuleIds)
		}

		if sharedKey := item["shared_key"]; sharedKey != "" {
			v.Properties.SharedKey = pointer.To(sharedKey.(string))
		}
		result = append(result, v)
	}

	return &result
}

func expandVpnGatewayConnectionIpSecPolicies(input []any) *[]virtualwans.IPsecPolicy {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.IPsecPolicy, 0)
	for _, itemRaw := range input {
		item := itemRaw.(map[string]any)
		result = append(result, virtualwans.IPsecPolicy{
			SaLifeTimeSeconds:   int64(item["sa_lifetime_sec"].(int)),
			SaDataSizeKilobytes: int64(item["sa_data_size_kb"].(int)),
			IPsecEncryption:     virtualwans.IPsecEncryption(item["encryption_algorithm"].(string)),
			IPsecIntegrity:      virtualwans.IPsecIntegrity(item["integrity_algorithm"].(string)),
			IkeEncryption:       virtualwans.IkeEncryption(item["ike_encryption_algorithm"].(string)),
			IkeIntegrity:        virtualwans.IkeIntegrity(item["ike_integrity_algorithm"].(string)),
			DhGroup:             virtualwans.DhGroup(item["dh_group"].(string)),
			PfsGroup:            virtualwans.PfsGroup(item["pfs_group"].(string)),
		})
	}

	return &result
}

func expandVpnGatewayConnectionRoutingConfiguration(input []any) *virtualwans.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	output := &virtualwans.RoutingConfiguration{
		AssociatedRouteTable: &virtualwans.SubResource{
			Id: pointer.To(raw["associated_route_table"].(string)),
		},
	}

	if inboundRouteMapId := raw["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		output.InboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := raw["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		output.OutboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	if v := raw["propagated_route_table"].([]any); len(v) != 0 {
		output.PropagatedRouteTables = expandVpnGatewayConnectionPropagatedRouteTable(v)
	}

	return output
}

func expandVpnGatewayConnectionTrafficSelectorPolicy(input []any) *[]virtualwans.TrafficSelectorPolicy {
	results := make([]virtualwans.TrafficSelectorPolicy, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, virtualwans.TrafficSelectorPolicy{
			LocalAddressRanges:  pointer.From(pluginsdk.ExpandStringSlice(v["local_address_ranges"].(*pluginsdk.Set).List())),
			RemoteAddressRanges: pointer.From(pluginsdk.ExpandStringSlice(v["remote_address_ranges"].(*pluginsdk.Set).List())),
		})
	}

	return &results
}

func expandVpnGatewayConnectionPropagatedRouteTable(input []any) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return &virtualwans.PropagatedRouteTable{}
	}

	v := input[0].(map[string]any)

	routeTableIds := make([]virtualwans.SubResource, 0)
	for _, val := range v["route_table_ids"].([]any) {
		routeTableIds = append(routeTableIds, virtualwans.SubResource{
			Id: pointer.To(val.(string)),
		})
	}

	result := virtualwans.PropagatedRouteTable{
		Ids: pointer.To(routeTableIds),
	}
	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = pluginsdk.ExpandStringSlice(labels)
	}
	return &result
}

func expandVpnGatewayConnectionNatRuleIds(input []any) *[]virtualwans.SubResource {
	results := make([]virtualwans.SubResource, 0)

	for _, item := range input {
		results = append(results, virtualwans.SubResource{
			Id: pointer.To(item.(string)),
		})
	}

	return &results
}

func expandVpnGatewayConnectionCustomBgpAddresses(input []any) *[]virtualwans.GatewayCustomBgpIPAddressIPConfiguration {
	results := make([]virtualwans.GatewayCustomBgpIPAddressIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, virtualwans.GatewayCustomBgpIPAddressIPConfiguration{
			CustomBgpIPAddress: v["ip_address"].(string),
			IPConfigurationId:  v["ip_configuration_id"].(string),
		})
	}

	return &results
}
