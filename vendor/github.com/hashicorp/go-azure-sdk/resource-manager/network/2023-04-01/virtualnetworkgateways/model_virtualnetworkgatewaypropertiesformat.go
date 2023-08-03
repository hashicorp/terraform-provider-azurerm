package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayPropertiesFormat struct {
	ActiveActive                      *bool                                   `json:"activeActive,omitempty"`
	AdminState                        *AdminState                             `json:"adminState,omitempty"`
	AllowRemoteVnetTraffic            *bool                                   `json:"allowRemoteVnetTraffic,omitempty"`
	AllowVirtualWanTraffic            *bool                                   `json:"allowVirtualWanTraffic,omitempty"`
	BgpSettings                       *BgpSettings                            `json:"bgpSettings,omitempty"`
	CustomRoutes                      *AddressSpace                           `json:"customRoutes,omitempty"`
	DisableIPSecReplayProtection      *bool                                   `json:"disableIPSecReplayProtection,omitempty"`
	EnableBgp                         *bool                                   `json:"enableBgp,omitempty"`
	EnableBgpRouteTranslationForNat   *bool                                   `json:"enableBgpRouteTranslationForNat,omitempty"`
	EnableDnsForwarding               *bool                                   `json:"enableDnsForwarding,omitempty"`
	EnablePrivateIPAddress            *bool                                   `json:"enablePrivateIpAddress,omitempty"`
	GatewayDefaultSite                *SubResource                            `json:"gatewayDefaultSite,omitempty"`
	GatewayType                       *VirtualNetworkGatewayType              `json:"gatewayType,omitempty"`
	IPConfigurations                  *[]VirtualNetworkGatewayIPConfiguration `json:"ipConfigurations,omitempty"`
	InboundDnsForwardingEndpoint      *string                                 `json:"inboundDnsForwardingEndpoint,omitempty"`
	NatRules                          *[]VirtualNetworkGatewayNatRule         `json:"natRules,omitempty"`
	ProvisioningState                 *ProvisioningState                      `json:"provisioningState,omitempty"`
	ResourceGuid                      *string                                 `json:"resourceGuid,omitempty"`
	Sku                               *VirtualNetworkGatewaySku               `json:"sku,omitempty"`
	VNetExtendedLocationResourceId    *string                                 `json:"vNetExtendedLocationResourceId,omitempty"`
	VirtualNetworkGatewayPolicyGroups *[]VirtualNetworkGatewayPolicyGroup     `json:"virtualNetworkGatewayPolicyGroups,omitempty"`
	VpnClientConfiguration            *VpnClientConfiguration                 `json:"vpnClientConfiguration,omitempty"`
	VpnGatewayGeneration              *VpnGatewayGeneration                   `json:"vpnGatewayGeneration,omitempty"`
	VpnType                           *VpnType                                `json:"vpnType,omitempty"`
}
