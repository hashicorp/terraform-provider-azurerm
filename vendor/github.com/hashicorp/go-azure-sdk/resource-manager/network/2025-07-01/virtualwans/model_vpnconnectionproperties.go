package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnConnectionProperties struct {
	ConnectionBandwidth            *int64                                   `json:"connectionBandwidth,omitempty"`
	ConnectionStatus               *VpnConnectionStatus                     `json:"connectionStatus,omitempty"`
	DpdTimeoutSeconds              *int64                                   `json:"dpdTimeoutSeconds,omitempty"`
	EgressBytesTransferred         *int64                                   `json:"egressBytesTransferred,omitempty"`
	EnableBgp                      *bool                                    `json:"enableBgp,omitempty"`
	EnableInternetSecurity         *bool                                    `json:"enableInternetSecurity,omitempty"`
	EnableRateLimiting             *bool                                    `json:"enableRateLimiting,omitempty"`
	IPsecPolicies                  *[]IPsecPolicy                           `json:"ipsecPolicies,omitempty"`
	IngressBytesTransferred        *int64                                   `json:"ingressBytesTransferred,omitempty"`
	ProvisioningState              *ProvisioningState                       `json:"provisioningState,omitempty"`
	RemoteVpnSite                  *SubResource                             `json:"remoteVpnSite,omitempty"`
	RoutingConfiguration           *RoutingConfiguration                    `json:"routingConfiguration,omitempty"`
	RoutingWeight                  *int64                                   `json:"routingWeight,omitempty"`
	SharedKey                      *string                                  `json:"sharedKey,omitempty"`
	TrafficSelectorPolicies        *[]TrafficSelectorPolicy                 `json:"trafficSelectorPolicies,omitempty"`
	UseLocalAzureIPAddress         *bool                                    `json:"useLocalAzureIpAddress,omitempty"`
	UsePolicyBasedTrafficSelectors *bool                                    `json:"usePolicyBasedTrafficSelectors,omitempty"`
	VpnConnectionProtocolType      *VirtualNetworkGatewayConnectionProtocol `json:"vpnConnectionProtocolType,omitempty"`
	VpnLinkConnections             *[]VpnSiteLinkConnection                 `json:"vpnLinkConnections,omitempty"`
}
