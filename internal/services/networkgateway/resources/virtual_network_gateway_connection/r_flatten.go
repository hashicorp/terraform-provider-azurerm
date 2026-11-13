// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
)

func flattenVirtualNetworkGatewayConnectionIpsecPolicies(ipsecPolicies *[]virtualnetworkgatewayconnections.IPsecPolicy) []any {
	schemaIpsecPolicies := make([]any, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]any)

			schemaIpsecPolicy["dh_group"] = string(ipsecPolicy.DhGroup)
			schemaIpsecPolicy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
			schemaIpsecPolicy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
			schemaIpsecPolicy["ipsec_encryption"] = string(ipsecPolicy.IPsecEncryption)
			schemaIpsecPolicy["ipsec_integrity"] = string(ipsecPolicy.IPsecIntegrity)
			schemaIpsecPolicy["pfs_group"] = string(ipsecPolicy.PfsGroup)
			schemaIpsecPolicy["sa_datasize"] = int(ipsecPolicy.SaDataSizeKilobytes)
			schemaIpsecPolicy["sa_lifetime"] = int(ipsecPolicy.SaLifeTimeSeconds)

			schemaIpsecPolicies = append(schemaIpsecPolicies, schemaIpsecPolicy)
		}
	}

	return schemaIpsecPolicies
}

func flattenGatewayCustomBgpIPAddresses(gatewayCustomBgpIPAddresses *[]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration) any {
	customBgpIpAddresses := make([]any, 0)
	if gatewayCustomBgpIPAddresses == nil || len(*gatewayCustomBgpIPAddresses) == 0 {
		return customBgpIpAddresses
	}

	customBgpIpAddress := map[string]any{}
	for k, v := range *gatewayCustomBgpIPAddresses {
		if k == 0 && v.CustomBgpIPAddress != "" {
			customBgpIpAddress["primary"] = v.CustomBgpIPAddress
		} else if k == 1 && v.CustomBgpIPAddress != "" {
			customBgpIpAddress["secondary"] = v.CustomBgpIPAddress
		}
	}

	return append(customBgpIpAddresses, customBgpIpAddress)
}

func flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(trafficSelectorPolicies *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy) []any {
	schemaTrafficSelectorPolicies := make([]any, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]any{
				"local_address_cidrs":  trafficSelectorPolicy.LocalAddressRanges,
				"remote_address_cidrs": trafficSelectorPolicy.RemoteAddressRanges,
			})
		}
	}

	return schemaTrafficSelectorPolicies
}

func flattenVirtualNetworkGatewayConnectionNatRuleIds(input *[]virtualnetworkgatewayconnections.SubResource) []any {
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
