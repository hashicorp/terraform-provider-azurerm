// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualNetworkGatewayConnectionIpsecPolicies(schemaIpsecPolicies []any) *[]virtualnetworkgatewayconnections.IPsecPolicy {
	ipsecPolicies := make([]virtualnetworkgatewayconnections.IPsecPolicy, 0, len(schemaIpsecPolicies))

	for _, d := range schemaIpsecPolicies {
		schemaIpsecPolicy := d.(map[string]any)
		ipsecPolicy := &virtualnetworkgatewayconnections.IPsecPolicy{}

		if dhGroup, ok := schemaIpsecPolicy["dh_group"].(string); ok && dhGroup != "" {
			ipsecPolicy.DhGroup = virtualnetworkgatewayconnections.DhGroup(dhGroup)
		}

		if ikeEncryption, ok := schemaIpsecPolicy["ike_encryption"].(string); ok && ikeEncryption != "" {
			ipsecPolicy.IkeEncryption = virtualnetworkgatewayconnections.IkeEncryption(ikeEncryption)
		}

		if ikeIntegrity, ok := schemaIpsecPolicy["ike_integrity"].(string); ok && ikeIntegrity != "" {
			ipsecPolicy.IkeIntegrity = virtualnetworkgatewayconnections.IkeIntegrity(ikeIntegrity)
		}

		if ipsecEncryption, ok := schemaIpsecPolicy["ipsec_encryption"].(string); ok && ipsecEncryption != "" {
			ipsecPolicy.IPsecEncryption = virtualnetworkgatewayconnections.IPsecEncryption(ipsecEncryption)
		}

		if ipsecIntegrity, ok := schemaIpsecPolicy["ipsec_integrity"].(string); ok && ipsecIntegrity != "" {
			ipsecPolicy.IPsecIntegrity = virtualnetworkgatewayconnections.IPsecIntegrity(ipsecIntegrity)
		}

		if pfsGroup, ok := schemaIpsecPolicy["pfs_group"].(string); ok && pfsGroup != "" {
			ipsecPolicy.PfsGroup = virtualnetworkgatewayconnections.PfsGroup(pfsGroup)
		}

		if v, ok := schemaIpsecPolicy["sa_datasize"].(int); ok {
			ipsecPolicy.SaDataSizeKilobytes = int64(v)
		}

		if v, ok := schemaIpsecPolicy["sa_lifetime"].(int); ok {
			ipsecPolicy.SaLifeTimeSeconds = int64(v)
		}

		ipsecPolicies = append(ipsecPolicies, *ipsecPolicy)
	}

	return &ipsecPolicies
}

func expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(schemaTrafficSelectorPolicies []any) *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy {
	trafficSelectorPolicies := make([]virtualnetworkgatewayconnections.TrafficSelectorPolicy, 0, len(schemaTrafficSelectorPolicies))

	for _, d := range schemaTrafficSelectorPolicies {
		schemaTrafficSelectorPolicy := d.(map[string]any)
		trafficSelectorPolicy := &virtualnetworkgatewayconnections.TrafficSelectorPolicy{}
		if localAddressRanges, ok := schemaTrafficSelectorPolicy["local_address_cidrs"].([]any); ok {
			trafficSelectorPolicy.LocalAddressRanges = pointer.From(pluginsdk.ExpandStringSlice(localAddressRanges))
		}
		if remoteAddressRanges, ok := schemaTrafficSelectorPolicy["remote_address_cidrs"].([]any); ok {
			trafficSelectorPolicy.RemoteAddressRanges = pointer.From(pluginsdk.ExpandStringSlice(remoteAddressRanges))
		}

		trafficSelectorPolicies = append(trafficSelectorPolicies, *trafficSelectorPolicy)
	}

	return &trafficSelectorPolicies
}

func expandGatewayCustomBgpIPAddresses(d *pluginsdk.ResourceData, bgpPeeringAddresses *[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress) (*[]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration, error) {
	customBgpIpAddresses := make([]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration, 0)

	bgpAddresses := d.Get("custom_bgp_addresses").([]any)
	if len(bgpAddresses) == 0 {
		return &customBgpIpAddresses, nil
	}

	bgAs := bgpAddresses[0].(map[string]any)
	primaryAddress := bgAs["primary"].(string)
	secondaryAddress := bgAs["secondary"].(string)

	var primaryIpConfiguration string
	var secondaryIpConfiguration string

	if bgpPeeringAddresses == nil {
		return &customBgpIpAddresses, fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway`, addresses returned nil")
	}

	for _, address := range *bgpPeeringAddresses {
		if address.CustomBgpIPAddresses == nil {
			continue
		}

		for _, ip := range *address.CustomBgpIPAddresses {
			if address.IPconfigurationId == nil {
				continue
			}

			switch ip {
			case primaryAddress:
				primaryIpConfiguration = *address.IPconfigurationId
			case secondaryAddress:
				secondaryIpConfiguration = *address.IPconfigurationId
			}
		}
	}

	if len(primaryIpConfiguration) == 0 || (secondaryAddress != "" && len(secondaryIpConfiguration) == 0) {
		return &customBgpIpAddresses, fmt.Errorf("primary or secondary address not found at `virtual_network_gateway` configuration `bgp_settings` `peering_addresses`")
	}

	customBgpIpAddresses = append(customBgpIpAddresses, virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration{
		IPConfigurationId:  primaryIpConfiguration,
		CustomBgpIPAddress: primaryAddress,
	})

	if secondaryAddress != "" {
		customBgpIpAddresses = append(customBgpIpAddresses, virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration{
			IPConfigurationId:  secondaryIpConfiguration,
			CustomBgpIPAddress: secondaryAddress,
		})
	}

	return &customBgpIpAddresses, nil
}

func expandVirtualNetworkGatewayConnectionNatRuleIds(input []any) *[]virtualnetworkgatewayconnections.SubResource {
	results := make([]virtualnetworkgatewayconnections.SubResource, 0)

	for _, item := range input {
		results = append(results, virtualnetworkgatewayconnections.SubResource{
			Id: pointer.To(item.(string)),
		})
	}

	return &results
}

func getVirtualNetworkGatewayConnectionProperties(d *pluginsdk.ResourceData, virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway) (*virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionType(d.Get("type").(string))
	connectionMode := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionMode(d.Get("connection_mode").(string))

	props := &virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType:                 connectionType,
		ConnectionMode:                 pointer.To(connectionMode),
		EnableBgp:                      pointer.To(d.Get("bgp_enabled").(bool)),
		EnablePrivateLinkFastPath:      pointer.To(d.Get("private_link_fast_path_enabled").(bool)),
		ExpressRouteGatewayBypass:      pointer.To(d.Get("express_route_gateway_bypass").(bool)),
		UsePolicyBasedTrafficSelectors: pointer.To(d.Get("use_policy_based_traffic_selectors").(bool)),
	}

	if virtualNetworkGateway.Name != nil && virtualNetworkGateway.Id != nil {
		props.VirtualNetworkGateway1 = virtualnetworkgatewayconnections.VirtualNetworkGateway{
			Id:   virtualNetworkGateway.Id,
			Name: virtualNetworkGateway.Name,
			Properties: virtualnetworkgatewayconnections.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]virtualnetworkgatewayconnections.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		props.AuthorizationKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("dpd_timeout_seconds"); ok {
		props.DpdTimeoutSeconds = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("express_route_circuit_id"); ok {
		props.Peer = &virtualnetworkgatewayconnections.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("egress_nat_rule_ids"); ok {
		props.EgressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("ingress_nat_rule_ids"); ok {
		props.IngressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("peer_virtual_network_gateway_id"); ok {
		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(v.(string))
		if err != nil {
			return nil, err
		}
		props.VirtualNetworkGateway2 = &virtualnetworkgatewayconnections.VirtualNetworkGateway{
			Id:   pointer.To(gwid.ID()),
			Name: &gwid.VirtualNetworkGatewayName,
			Properties: virtualnetworkgatewayconnections.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]virtualnetworkgatewayconnections.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("local_azure_ip_address_enabled"); ok {
		props.UseLocalAzureIPAddress = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("local_network_gateway_id"); ok {
		localNetworkGatewayId := v.(string)
		name, err := localNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("getting LocalNetworkGateway Name and Group: %+v", err)
		}

		props.LocalNetworkGateway2 = &virtualnetworkgatewayconnections.LocalNetworkGateway{
			Id:   &localNetworkGatewayId,
			Name: &name,
			Properties: virtualnetworkgatewayconnections.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &virtualnetworkgatewayconnections.AddressSpace{},
			},
		}
	}

	if v, ok := d.GetOk("routing_weight"); ok {
		props.RoutingWeight = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("shared_key"); ok {
		props.SharedKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("connection_protocol"); ok {
		connectionProtocol := v.(string)
		props.ConnectionProtocol = pointer.ToEnum[virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionProtocol](connectionProtocol)
	}

	if v, ok := d.GetOk("traffic_selector_policy"); ok {
		props.TrafficSelectorPolicies = expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(v.([]any))
	}

	if v, ok := d.GetOk("ipsec_policy"); ok {
		props.IPsecPolicies = expandVirtualNetworkGatewayConnectionIpsecPolicies(v.([]any))
	}

	if pointer.From(props.EnableBgp) {
		if _, ok := d.GetOk("custom_bgp_addresses"); ok {
			if virtualNetworkGateway.Properties.BgpSettings == nil || virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses == nil {
				return nil, fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway` %s (%s) failed: get nil", *virtualNetworkGateway.Name, *virtualNetworkGateway.Id)
			}

			gatewayCustomBgpIPAddresses, err := expandGatewayCustomBgpIPAddresses(d, virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses)
			if err != nil {
				return nil, err
			}

			props.GatewayCustomBgpIPAddresses = gatewayCustomBgpIPAddresses
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute {
		if props.Peer == nil || props.Peer.Id == nil {
			return nil, fmt.Errorf("`express_route_circuit_id` must be specified when `type` is set to `ExpressRoute`")
		}
		if d.Get("private_link_fast_path_enabled").(bool) && !d.Get("express_route_gateway_bypass").(bool) {
			return nil, fmt.Errorf("`express_route_gateway_bypass` must be enabled when `private_link_fast_path_enabled` is set to `true`")
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec {
		if props.LocalNetworkGateway2 == nil || props.LocalNetworkGateway2.Id == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` must be specified when `type` is set to `IPsec`")
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet {
		if props.VirtualNetworkGateway2 == nil || props.VirtualNetworkGateway2.Id == nil {
			return nil, fmt.Errorf("`peer_virtual_network_gateway_id` must be specified when `type` is set to `Vnet2Vnet`")
		}
	}

	if props.GatewayCustomBgpIPAddresses != nil && props.ConnectionType != virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `type` is set to `IPsec`")
	}

	if props.GatewayCustomBgpIPAddresses != nil && virtualNetworkGateway.Properties.ActiveActive == pointer.To(false) {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `azurerm_virtual_network_gateway` `active_active` is set enabled`")
	}

	return props, nil
}
