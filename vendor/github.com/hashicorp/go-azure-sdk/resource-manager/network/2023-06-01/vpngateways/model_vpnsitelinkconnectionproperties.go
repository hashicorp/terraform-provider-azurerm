package vpngateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnSiteLinkConnectionProperties struct {
	ConnectionBandwidth            *int64                                      `json:"connectionBandwidth,omitempty"`
	ConnectionStatus               *VpnConnectionStatus                        `json:"connectionStatus,omitempty"`
	EgressBytesTransferred         *int64                                      `json:"egressBytesTransferred,omitempty"`
	EgressNatRules                 *[]SubResource                              `json:"egressNatRules,omitempty"`
	EnableBgp                      *bool                                       `json:"enableBgp,omitempty"`
	EnableRateLimiting             *bool                                       `json:"enableRateLimiting,omitempty"`
	IPsecPolicies                  *[]IPsecPolicy                              `json:"ipsecPolicies,omitempty"`
	IngressBytesTransferred        *int64                                      `json:"ingressBytesTransferred,omitempty"`
	IngressNatRules                *[]SubResource                              `json:"ingressNatRules,omitempty"`
	ProvisioningState              *ProvisioningState                          `json:"provisioningState,omitempty"`
	RoutingWeight                  *int64                                      `json:"routingWeight,omitempty"`
	SharedKey                      *string                                     `json:"sharedKey,omitempty"`
	UseLocalAzureIPAddress         *bool                                       `json:"useLocalAzureIpAddress,omitempty"`
	UsePolicyBasedTrafficSelectors *bool                                       `json:"usePolicyBasedTrafficSelectors,omitempty"`
	VpnConnectionProtocolType      *VirtualNetworkGatewayConnectionProtocol    `json:"vpnConnectionProtocolType,omitempty"`
	VpnGatewayCustomBgpAddresses   *[]GatewayCustomBgpIPAddressIPConfiguration `json:"vpnGatewayCustomBgpAddresses,omitempty"`
	VpnLinkConnectionMode          *VpnLinkConnectionMode                      `json:"vpnLinkConnectionMode,omitempty"`
	VpnSiteLink                    *SubResource                                `json:"vpnSiteLink,omitempty"`
}
