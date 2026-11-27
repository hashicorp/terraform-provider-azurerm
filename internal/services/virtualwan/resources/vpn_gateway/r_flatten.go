// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenVPNGatewayBGPSettings(input *virtualwans.BgpSettings) []any {
	if input == nil {
		return []any{}
	}

	asn := 0
	if input.Asn != nil {
		asn = int(*input.Asn)
	}

	peerWeight := 0
	if input.PeerWeight != nil {
		peerWeight = int(*input.PeerWeight)
	}

	var instance0BgpPeeringAddress, instance1BgpPeeringAddress []any
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 0 {
		instance0BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[0])
	}
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 1 {
		instance1BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[1])
	}

	return []any{
		map[string]any{
			"asn":                            asn,
			"bgp_peering_address":            pointer.From(input.BgpPeeringAddress),
			"instance_0_bgp_peering_address": instance0BgpPeeringAddress,
			"instance_1_bgp_peering_address": instance1BgpPeeringAddress,
			"peer_weight":                    peerWeight,
		},
	}
}

func flattenVPNGatewayIPConfigurationBgpPeeringAddress(input virtualwans.IPConfigurationBgpPeeringAddress) []any {
	return []any{
		map[string]any{
			"ip_configuration_id": pointer.From(input.IPconfigurationId),
			"custom_ips":          pluginsdk.FlattenSlice(input.CustomBgpIPAddresses),
			"default_ips":         pluginsdk.FlattenSlice(input.DefaultBgpIPAddresses),
			"tunnel_ips":          pluginsdk.FlattenSlice(input.TunnelIPAddresses),
		},
	}
}

func flattenVPNGatewayIpConfiguration(input *[]virtualwans.VpnGatewayIPConfiguration) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, map[string]any{
			"id":                 pointer.From(item.Id),
			"private_ip_address": pointer.From(item.PrivateIPAddress),
			"public_ip_address":  pointer.From(item.PublicIPAddress),
		})
	}

	return result
}
