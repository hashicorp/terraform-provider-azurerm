// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenVirtualNetworkGatewayBgpSettings(settings *virtualnetworkgateways.BgpSettings) ([]any, error) {
	output := make([]any, 0)

	if settings != nil {
		flat := make(map[string]any)

		if asn := settings.Asn; asn != nil {
			flat["asn"] = int(*asn)
		}
		if weight := settings.PeerWeight; weight != nil {
			flat["peer_weight"] = int(*weight)
		}

		var err error
		flat["peering_addresses"], err = flattenVirtualNetworkGatewayBgpPeeringAddresses(settings.BgpPeeringAddresses)
		if err != nil {
			return nil, err
		}

		output = append(output, flat)
	}

	return output, nil
}

func flattenVirtualNetworkGatewayBgpPeeringAddresses(input *[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress) (any, error) {
	if input == nil {
		return []any{}, nil
	}

	output := make([]any, 0)

	for _, e := range *input {
		var ipConfigName string
		if e.IPconfigurationId != nil {
			id, err := parse.VirtualNetworkGatewayIpConfigurationIDInsensitively(*e.IPconfigurationId)
			if err != nil {
				return nil, err
			}
			ipConfigName = id.IpConfigurationName
		}

		output = append(output, map[string]any{
			"ip_configuration_name": ipConfigName,
			"apipa_addresses":       pluginsdk.FlattenSlice(e.CustomBgpIPAddresses),
			"default_addresses":     pluginsdk.FlattenSlice(e.DefaultBgpIPAddresses),
			"tunnel_ip_addresses":   pluginsdk.FlattenSlice(e.TunnelIPAddresses),
		})
	}

	return output, nil
}

func flattenVirtualNetworkGatewayIPConfigurations(ipConfigs *[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration, gatewayType virtualnetworkgateways.VirtualNetworkGatewayType) []any {
	flat := make([]any, 0)

	if ipConfigs != nil {
		for _, cfg := range *ipConfigs {
			props := cfg.Properties
			v := make(map[string]any)

			if name := cfg.Name; name != nil {
				v["name"] = *name
			}
			v["private_ip_address_allocation"] = pointer.FromEnum(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				if id := subnet.Id; id != nil {
					v["subnet_id"] = *id
				}
			}

			// Do not include public_ip_address_id for ExpressRoute gateways
			if gatewayType != virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute {
				if pip := props.PublicIPAddress; pip != nil {
					if id := pip.Id; id != nil {
						v["public_ip_address_id"] = *id
					}
				}
			}

			flat = append(flat, v)
		}
	}

	return flat
}

func flattenVirtualNetworkGatewayVpnClientConfig(cfg *virtualnetworkgateways.VpnClientConfiguration) ([]any, error) {
	if cfg == nil {
		return []any{}, nil
	}
	flat := map[string]any{
		"ipsec_policy":  flattenVirtualNetworkGatewayIPSecPolicies(cfg.VpnClientIPsecPolicies),
		"radius_server": flattenVirtualNetworkGatewayRadiusServers(cfg.RadiusServers),
	}

	connection, err := flattenVirtualNetworkGatewayClientConnections(cfg.VngClientConnectionConfigurations)
	if err != nil {
		return nil, err
	}
	flat["virtual_network_gateway_client_connection"] = connection

	if pool := cfg.VpnClientAddressPool; pool != nil {
		flat["address_space"] = pluginsdk.FlattenSlice(pool.AddressPrefixes)
	} else {
		flat["address_space"] = []any{}
	}

	if v := cfg.AadTenant; v != nil {
		flat["aad_tenant"] = *v
	}

	if v := cfg.AadAudience; v != nil {
		flat["aad_audience"] = *v
	}

	if v := cfg.AadIssuer; v != nil {
		flat["aad_issuer"] = *v
	}

	rootCerts := make([]any, 0)
	if certs := cfg.VpnClientRootCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]any{
				"name":             *cert.Name,
				"public_cert_data": cert.Properties.PublicCertData,
			}
			rootCerts = append(rootCerts, v)
		}
	}
	flat["root_certificate"] = pluginsdk.NewSet(hashVirtualNetworkGatewayRootCert, rootCerts)

	revokedCerts := make([]any, 0)
	if certs := cfg.VpnClientRevokedCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]any{
				"name":       *cert.Name,
				"thumbprint": *cert.Properties.Thumbprint,
			}
			revokedCerts = append(revokedCerts, v)
		}
	}
	flat["revoked_certificate"] = pluginsdk.NewSet(hashVirtualNetworkGatewayRevokedCert, revokedCerts)

	vpnClientProtocols := &pluginsdk.Set{F: pluginsdk.HashString}
	if vpnProtocols := cfg.VpnClientProtocols; vpnProtocols != nil {
		for _, protocol := range *vpnProtocols {
			vpnClientProtocols.Add(string(protocol))
		}
	}
	flat["vpn_client_protocols"] = vpnClientProtocols

	vpnAuthTypes := &pluginsdk.Set{F: pluginsdk.HashString}
	if authTypes := cfg.VpnAuthenticationTypes; authTypes != nil {
		for _, authType := range *authTypes {
			vpnAuthTypes.Add(string(authType))
		}
	}
	flat["vpn_auth_types"] = vpnAuthTypes

	if v := cfg.RadiusServerAddress; v != nil {
		flat["radius_server_address"] = *v
	}

	if v := cfg.RadiusServerSecret; v != nil {
		flat["radius_server_secret"] = *v
	}

	return []any{flat}, nil
}

func flattenVirtualNetworkGatewayAddressSpace(input *virtualnetworkgateways.AddressSpace) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"address_prefixes": pluginsdk.FlattenSlice(input.AddressPrefixes),
		},
	}
}

func flattenVirtualNetworkGatewayRadiusServers(input *[]virtualnetworkgateways.RadiusServer) []any {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"address": item.RadiusServerAddress,
			"secret":  pointer.From(item.RadiusServerSecret),
			"score":   pointer.From(item.RadiusServerScore),
		})
	}

	return results
}

func flattenVirtualNetworkGatewayIPSecPolicies(input *[]virtualnetworkgateways.IPsecPolicy) []any {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"dh_group":                  string(item.DhGroup),
			"ipsec_encryption":          string(item.IPsecEncryption),
			"ipsec_integrity":           string(item.IPsecIntegrity),
			"ike_encryption":            string(item.IkeEncryption),
			"ike_integrity":             string(item.IkeIntegrity),
			"pfs_group":                 string(item.PfsGroup),
			"sa_data_size_in_kilobytes": item.SaDataSizeKilobytes,
			"sa_lifetime_in_seconds":    item.SaLifeTimeSeconds,
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicyGroups(input *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup) []any {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		var isDefault bool
		var policyMember any
		var priority int64
		if props := item.Properties; props != nil {
			isDefault = props.IsDefault
			policyMember = flattenVirtualNetworkGatewayPolicy(props.PolicyMembers)
			priority = props.Priority
		}
		results = append(results, map[string]any{
			"name":          pointer.From(item.Name),
			"is_default":    isDefault,
			"policy_member": policyMember,
			"priority":      priority,
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicy(input []virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember) []any {
	results := make([]any, 0)
	if len(input) == 0 {
		return results
	}

	for _, item := range input {
		results = append(results, map[string]any{
			"name":  pointer.From(item.Name),
			"type":  pointer.FromEnum(item.AttributeType),
			"value": pointer.From(item.AttributeValue),
		})
	}
	return results
}

func flattenVirtualNetworkGatewayClientConnections(input *[]virtualnetworkgateways.VngClientConnectionConfiguration) ([]any, error) {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, item := range *input {
		var err error
		addressPrefixes := make([]string, 0)
		policyGroupNames := make([]string, 0)
		if props := item.Properties; props != nil {
			addressPrefixes = pointer.From(props.VpnClientAddressPool.AddressPrefixes)
			policyGroupNames, err = flattenVirtualNetworkGatewayPolicyGroupNames(props.VirtualNetworkGatewayPolicyGroups)
			if err != nil {
				return nil, err
			}
		}
		result := map[string]any{
			"name":             pointer.From(item.Name),
			"address_prefixes": addressPrefixes,
		}

		result["policy_group_names"] = policyGroupNames

		results = append(results, result)
	}

	return results, nil
}

func flattenVirtualNetworkGatewayPolicyGroupNames(input []virtualnetworkgateways.SubResource) ([]string, error) {
	results := make([]string, 0)
	if len(input) == 0 {
		return results, nil
	}

	for _, item := range input {
		policyGroupId, err := parse.VirtualNetworkGatewayPolicyGroupID(*item.Id)
		if err != nil {
			return nil, err
		}
		results = append(results, policyGroupId.Name)
	}

	return results, nil
}

func flattenVirtualNetworkGatewayAutoScaleConfiguration(input *virtualnetworkgateways.VirtualNetworkGatewayAutoScaleConfiguration) (any, any) {
	if input == nil || input.Bounds == nil {
		return nil, nil
	}

	return input.Bounds.Min, input.Bounds.Max
}
