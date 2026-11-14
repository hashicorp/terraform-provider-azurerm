// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualNetworkGatewayBgpSettings(id virtualnetworkgateways.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*virtualnetworkgateways.BgpSettings, error) {
	bgpSets := d.Get("bgp_settings").([]any)
	if len(bgpSets) == 0 {
		return nil, nil
	}

	bgp := bgpSets[0].(map[string]any)

	peeringAddresses, err := expandVirtualNetworkGatewayBgpPeeringAddresses(id, d.Get("ip_configuration").([]any), bgp["peering_addresses"].([]any))
	if err != nil {
		return nil, err
	}

	return &virtualnetworkgateways.BgpSettings{
		Asn:                 pointer.To(int64(bgp["asn"].(int))),
		PeerWeight:          pointer.To(int64(bgp["peer_weight"].(int))),
		BgpPeeringAddresses: peeringAddresses,
	}, nil
}

func expandVirtualNetworkGatewayBgpPeeringAddresses(id virtualnetworkgateways.VirtualNetworkGatewayId, ipConfig, input []any) (*[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]virtualnetworkgateways.IPConfigurationBgpPeeringAddress, 0)

	var existIpConfigName string
	if len(ipConfig) == 1 {
		existIpConfigName = ipConfig[0].(map[string]any)["name"].(string)
	}

	for _, e := range input {
		if e == nil {
			continue
		}
		b := e.(map[string]any)

		ipConfigName := b["ip_configuration_name"].(string)

		if ipConfigName == "" {
			// existIpConfigName is empty means there are more than one `ip_configuration` blocks, in which case users have to specify the
			// `ip_configuration_name` used in current `peering_addresses` setting.
			if existIpConfigName == "" {
				return nil, fmt.Errorf("`ip_configuration_name` has to be set in current `peering_addresses` block in case there are multiple `ip_configuration` blocks")
			}

			ipConfigName = existIpConfigName
		}

		// If there is an existing ip configuration name defined in the only `ip_configuration` block, and users explicitly set the `ip_configuration_name` in current
		// `peering_addresses` block, they should be the same name.
		if existIpConfigName != "" && ipConfigName != "" {
			if ipConfigName != existIpConfigName {
				return nil, fmt.Errorf("`ip_configuration.0.name` is not the same as `bgp_settings.0.peering_addresses.*.ip_configuration_name`")
			}
		}

		ipConfigId := parse.NewVirtualNetworkGatewayIpConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName, ipConfigName)
		result = append(result, virtualnetworkgateways.IPConfigurationBgpPeeringAddress{
			IPconfigurationId:    pointer.To(ipConfigId.ID()),
			CustomBgpIPAddresses: pluginsdk.ExpandStringSlice(b["apipa_addresses"].([]any)),
		})
	}

	return &result, nil
}

func expandVirtualNetworkGatewayIPConfigurations(d *pluginsdk.ResourceData) *[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration {
	configs := d.Get("ip_configuration").([]any)
	ipConfigs := make([]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration, 0, len(configs))

	for _, c := range configs {
		conf := c.(map[string]any)

		props := &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: pointer.ToEnum[virtualnetworkgateways.IPAllocationMethod](conf["private_ip_address_allocation"].(string)),
		}

		if subnetID := conf["subnet_id"].(string); subnetID != "" {
			props.Subnet = &virtualnetworkgateways.SubResource{
				Id: &subnetID,
			}
		}

		if publicIP := conf["public_ip_address_id"].(string); publicIP != "" {
			props.PublicIPAddress = &virtualnetworkgateways.SubResource{
				Id: &publicIP,
			}
		}

		ipConfig := virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
			Name:       pointer.To(conf["name"].(string)),
			Properties: props,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return &ipConfigs
}

func expandVirtualNetworkGatewayVpnClientConfig(d *pluginsdk.ResourceData, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *virtualnetworkgateways.VpnClientConfiguration {
	configSets := d.Get("vpn_client_configuration").([]any)
	if len(configSets) == 0 {
		// return nil will delete the existing vpn client configuration
		return nil
	}

	conf := configSets[0].(map[string]any)

	confAddresses := conf["address_space"].([]any)
	addresses := make([]string, 0, len(confAddresses))
	for _, addr := range confAddresses {
		addresses = append(addresses, addr.(string))
	}

	rootCertsConf := conf["root_certificate"].(*pluginsdk.Set).List()
	rootCerts := make([]virtualnetworkgateways.VpnClientRootCertificate, 0, len(rootCertsConf))
	for _, rootCertSet := range rootCertsConf {
		rootCert := rootCertSet.(map[string]any)
		r := virtualnetworkgateways.VpnClientRootCertificate{
			Name: pointer.To(rootCert["name"].(string)),
			Properties: virtualnetworkgateways.VpnClientRootCertificatePropertiesFormat{
				PublicCertData: rootCert["public_cert_data"].(string),
			},
		}
		rootCerts = append(rootCerts, r)
	}

	revokedCertsConf := conf["revoked_certificate"].(*pluginsdk.Set).List()
	revokedCerts := make([]virtualnetworkgateways.VpnClientRevokedCertificate, 0, len(revokedCertsConf))
	for _, revokedCertSet := range revokedCertsConf {
		revokedCert := revokedCertSet.(map[string]any)
		r := virtualnetworkgateways.VpnClientRevokedCertificate{
			Name: pointer.To(revokedCert["name"].(string)),
			Properties: &virtualnetworkgateways.VpnClientRevokedCertificatePropertiesFormat{
				Thumbprint: pointer.To(revokedCert["thumbprint"].(string)),
			},
		}
		revokedCerts = append(revokedCerts, r)
	}

	vpnClientProtocolsConf := conf["vpn_client_protocols"].(*pluginsdk.Set).List()
	vpnClientProtocols := make([]virtualnetworkgateways.VpnClientProtocol, 0, len(vpnClientProtocolsConf))
	for _, vpnClientProtocol := range vpnClientProtocolsConf {
		p := virtualnetworkgateways.VpnClientProtocol(vpnClientProtocol.(string))
		vpnClientProtocols = append(vpnClientProtocols, p)
	}

	vpnAuthTypesConf := conf["vpn_auth_types"].(*pluginsdk.Set).List()
	vpnAuthTypes := make([]virtualnetworkgateways.VpnAuthenticationType, 0, len(vpnAuthTypesConf))
	for _, vpnAuthType := range vpnAuthTypesConf {
		a := virtualnetworkgateways.VpnAuthenticationType(vpnAuthType.(string))
		vpnAuthTypes = append(vpnAuthTypes, a)
	}

	return &virtualnetworkgateways.VpnClientConfiguration{
		VpnClientAddressPool: &virtualnetworkgateways.AddressSpace{
			AddressPrefixes: &addresses,
		},
		AadTenant:                         pointer.To(conf["aad_tenant"].(string)),
		AadAudience:                       pointer.To(conf["aad_audience"].(string)),
		AadIssuer:                         pointer.To(conf["aad_issuer"].(string)),
		VngClientConnectionConfigurations: expandVirtualNetworkGatewayClientConnections(conf["virtual_network_gateway_client_connection"].([]any), vnetGatewayId),
		VpnClientIPsecPolicies:            expandVirtualNetworkGatewayIpsecPolicies(conf["ipsec_policy"].([]any)),
		VpnClientRootCertificates:         &rootCerts,
		VpnClientRevokedCertificates:      &revokedCerts,
		VpnClientProtocols:                &vpnClientProtocols,
		RadiusServers:                     expandVirtualNetworkGatewayRadiusServers(conf["radius_server"].([]any)),
		RadiusServerAddress:               pointer.To(conf["radius_server_address"].(string)),
		RadiusServerSecret:                pointer.To(conf["radius_server_secret"].(string)),
		VpnAuthenticationTypes:            &vpnAuthTypes,
	}
}

func expandVirtualNetworkGatewaySku(d *pluginsdk.ResourceData) *virtualnetworkgateways.VirtualNetworkGatewaySku {
	sku := d.Get("sku").(string)

	return &virtualnetworkgateways.VirtualNetworkGatewaySku{
		Name: pointer.ToEnum[virtualnetworkgateways.VirtualNetworkGatewaySkuName](sku),
		Tier: pointer.ToEnum[virtualnetworkgateways.VirtualNetworkGatewaySkuTier](sku),
	}
}

func expandVirtualNetworkGatewayAddressSpace(input []any) *virtualnetworkgateways.AddressSpace {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &virtualnetworkgateways.AddressSpace{
		AddressPrefixes: pluginsdk.ExpandStringSlice(v["address_prefixes"].(*pluginsdk.Set).List()),
	}
}

func expandVirtualNetworkGatewayIpsecPolicies(input []any) *[]virtualnetworkgateways.IPsecPolicy {
	results := make([]virtualnetworkgateways.IPsecPolicy, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, virtualnetworkgateways.IPsecPolicy{
			DhGroup:             virtualnetworkgateways.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       virtualnetworkgateways.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        virtualnetworkgateways.IkeIntegrity(v["ike_integrity"].(string)),
			IPsecEncryption:     virtualnetworkgateways.IPsecEncryption(v["ipsec_encryption"].(string)),
			IPsecIntegrity:      virtualnetworkgateways.IPsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            virtualnetworkgateways.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   int64(v["sa_lifetime_in_seconds"].(int)),
			SaDataSizeKilobytes: int64(v["sa_data_size_in_kilobytes"].(int)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayRadiusServers(input []any) *[]virtualnetworkgateways.RadiusServer {
	results := make([]virtualnetworkgateways.RadiusServer, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, virtualnetworkgateways.RadiusServer{
			RadiusServerAddress: v["address"].(string),
			RadiusServerScore:   pointer.To(int64(v["score"].(int))),
			RadiusServerSecret:  pointer.To(v["secret"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyGroups(input []any) *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup {
	results := make([]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyGroup := item.(map[string]any)

		results = append(results, virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup{
			Name: pointer.To(policyGroup["name"].(string)),
			Properties: &virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupProperties{
				IsDefault:     policyGroup["is_default"].(bool),
				PolicyMembers: *expandVirtualNetworkGatewayPolicyMembers(policyGroup["policy_member"].([]any)),
				Priority:      int64(policyGroup["priority"].(int)),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyMembers(input []any) *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember {
	results := make([]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyMember := item.(map[string]any)

		results = append(results, virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember{
			Name:           pointer.To(policyMember["name"].(string)),
			AttributeType:  pointer.ToEnum[virtualnetworkgateways.VpnPolicyMemberAttributeType](policyMember["type"].(string)),
			AttributeValue: pointer.To(policyMember["value"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayClientConnections(input []any, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *[]virtualnetworkgateways.VngClientConnectionConfiguration {
	results := make([]virtualnetworkgateways.VngClientConnectionConfiguration, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		vngClientConnectionConfiguration := item.(map[string]any)

		results = append(results, virtualnetworkgateways.VngClientConnectionConfiguration{
			Name: pointer.To(vngClientConnectionConfiguration["name"].(string)),
			Properties: &virtualnetworkgateways.VngClientConnectionConfigurationProperties{
				VpnClientAddressPool:              *expandVirtualNetworkGatewayAddressPool(vngClientConnectionConfiguration["address_prefixes"].([]any)),
				VirtualNetworkGatewayPolicyGroups: *expandVirtualNetworkGatewayPolicyGroupNames(vngClientConnectionConfiguration["policy_group_names"].([]any), vnetGatewayId),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayAddressPool(input []any) *virtualnetworkgateways.AddressSpace {
	if len(input) == 0 {
		return &virtualnetworkgateways.AddressSpace{}
	}

	addressPrefixes := make([]string, 0)
	for _, v := range input {
		addressPrefixes = append(addressPrefixes, v.(string))
	}

	return &virtualnetworkgateways.AddressSpace{
		AddressPrefixes: pointer.To(addressPrefixes),
	}
}

func expandVirtualNetworkGatewayPolicyGroupNames(input []any, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *[]virtualnetworkgateways.SubResource {
	results := make([]virtualnetworkgateways.SubResource, 0)
	if len(input) == 0 {
		return &results
	}

	for _, v := range input {
		policyGroupId := parse.NewVirtualNetworkGatewayPolicyGroupID(vnetGatewayId.SubscriptionId, vnetGatewayId.ResourceGroupName, vnetGatewayId.VirtualNetworkGatewayName, v.(string))
		result := virtualnetworkgateways.SubResource{
			Id: pointer.To(policyGroupId.ID()),
		}
		results = append(results, result)
	}

	return &results
}

func expandVirtualNetworkGatewayAutoScaleConfiguration(d *pluginsdk.ResourceData) *virtualnetworkgateways.VirtualNetworkGatewayAutoScaleConfiguration {
	minScaleUnit := d.Get("minimum_scale_unit").(int)
	maxScaleUnit := d.Get("maximum_scale_unit").(int)

	if minScaleUnit == 0 {
		return nil
	}

	return &virtualnetworkgateways.VirtualNetworkGatewayAutoScaleConfiguration{
		Bounds: &virtualnetworkgateways.VirtualNetworkGatewayAutoScaleBounds{
			Min: pointer.To(int64(minScaleUnit)),
			Max: pointer.To(int64(maxScaleUnit)),
		},
	}
}

func getVirtualNetworkGatewayProperties(id virtualnetworkgateways.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*virtualnetworkgateways.VirtualNetworkGatewayPropertiesFormat, error) {
	props := &virtualnetworkgateways.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:                     pointer.ToEnum[virtualnetworkgateways.VirtualNetworkGatewayType](d.Get("type").(string)),
		VpnType:                         pointer.ToEnum[virtualnetworkgateways.VpnType](d.Get("vpn_type").(string)),
		EnableBgp:                       pointer.To(d.Get("bgp_enabled").(bool)),
		EnablePrivateIPAddress:          pointer.To(d.Get("private_ip_address_enabled").(bool)),
		ActiveActive:                    pointer.To(d.Get("active_active").(bool)),
		EnableBgpRouteTranslationForNat: pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool)),
		DisableIPSecReplayProtection:    pointer.To(!d.Get("ip_sec_replay_protection_enabled").(bool)),
		AllowRemoteVnetTraffic:          pointer.To(d.Get("remote_vnet_traffic_enabled").(bool)),
		AllowVirtualWanTraffic:          pointer.To(d.Get("virtual_wan_traffic_enabled").(bool)),
		Sku:                             expandVirtualNetworkGatewaySku(d),
		IPConfigurations:                expandVirtualNetworkGatewayIPConfigurations(d),
		CustomRoutes:                    expandVirtualNetworkGatewayAddressSpace(d.Get("custom_route").([]any)),
	}

	if v, ok := d.GetOk("generation"); ok {
		props.VpnGatewayGeneration = pointer.ToEnum[virtualnetworkgateways.VpnGatewayGeneration](v.(string))
	}

	if v, ok := d.GetOk("dns_forwarding_enabled"); ok {
		props.EnableDnsForwarding = pointer.To(v.(bool))
	}

	if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
		props.GatewayDefaultSite = &virtualnetworkgateways.SubResource{
			Id: &gatewayDefaultSiteID,
		}
	}

	props.AutoScaleConfiguration = expandVirtualNetworkGatewayAutoScaleConfiguration(d)

	if v, ok := d.GetOk("policy_group"); ok {
		props.VirtualNetworkGatewayPolicyGroups = expandVirtualNetworkGatewayPolicyGroups(v.([]any))
	}

	if _, ok := d.GetOk("vpn_client_configuration"); ok {
		props.VpnClientConfiguration = expandVirtualNetworkGatewayVpnClientConfig(d, id)
	}

	if _, ok := d.GetOk("bgp_settings"); ok {
		bgpSettings, err := expandVirtualNetworkGatewayBgpSettings(id, d)
		if err != nil {
			return nil, err
		}
		props.BgpSettings = bgpSettings
	}

	gatewayType := pointer.From(props.GatewayType)
	vpnType := pointer.From(props.VpnType)
	vpnGatewayGeneration := pointer.From(props.VpnGatewayGeneration)
	skuName := pointer.FromEnum(props.Sku.Name)

	// Sku validation for policy-based VPN gateways
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypePolicyBased {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayPolicyBasedVpnSku()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of first geneneration
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypeRouteBased && vpnGatewayGeneration == virtualnetworkgateways.VpnGatewayGenerationGenerationOne {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of second geneneration
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypeRouteBased && vpnGatewayGeneration == virtualnetworkgateways.VpnGatewayGenerationGenerationTwo {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2()); !ok {
			return nil, err
		}
	}

	// Sku validation for ExpressRoute gateways
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayExpressRouteSku()); !ok {
			return nil, err
		}
	}

	return props, nil
}
