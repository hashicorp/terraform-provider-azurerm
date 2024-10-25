package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayConnectionPropertiesFormat struct {
	AuthorizationKey               *string                                     `json:"authorizationKey,omitempty"`
	ConnectionMode                 *VirtualNetworkGatewayConnectionMode        `json:"connectionMode,omitempty"`
	ConnectionProtocol             *VirtualNetworkGatewayConnectionProtocol    `json:"connectionProtocol,omitempty"`
	ConnectionStatus               *VirtualNetworkGatewayConnectionStatus      `json:"connectionStatus,omitempty"`
	ConnectionType                 VirtualNetworkGatewayConnectionType         `json:"connectionType"`
	DpdTimeoutSeconds              *int64                                      `json:"dpdTimeoutSeconds,omitempty"`
	EgressBytesTransferred         *int64                                      `json:"egressBytesTransferred,omitempty"`
	EgressNatRules                 *[]SubResource                              `json:"egressNatRules,omitempty"`
	EnableBgp                      *bool                                       `json:"enableBgp,omitempty"`
	EnablePrivateLinkFastPath      *bool                                       `json:"enablePrivateLinkFastPath,omitempty"`
	ExpressRouteGatewayBypass      *bool                                       `json:"expressRouteGatewayBypass,omitempty"`
	GatewayCustomBgpIPAddresses    *[]GatewayCustomBgpIPAddressIPConfiguration `json:"gatewayCustomBgpIpAddresses,omitempty"`
	IPsecPolicies                  *[]IPsecPolicy                              `json:"ipsecPolicies,omitempty"`
	IngressBytesTransferred        *int64                                      `json:"ingressBytesTransferred,omitempty"`
	IngressNatRules                *[]SubResource                              `json:"ingressNatRules,omitempty"`
	LocalNetworkGateway2           *LocalNetworkGateway                        `json:"localNetworkGateway2,omitempty"`
	Peer                           *SubResource                                `json:"peer,omitempty"`
	ProvisioningState              *ProvisioningState                          `json:"provisioningState,omitempty"`
	ResourceGuid                   *string                                     `json:"resourceGuid,omitempty"`
	RoutingWeight                  *int64                                      `json:"routingWeight,omitempty"`
	SharedKey                      *string                                     `json:"sharedKey,omitempty"`
	TrafficSelectorPolicies        *[]TrafficSelectorPolicy                    `json:"trafficSelectorPolicies,omitempty"`
	TunnelConnectionStatus         *[]TunnelConnectionHealth                   `json:"tunnelConnectionStatus,omitempty"`
	UseLocalAzureIPAddress         *bool                                       `json:"useLocalAzureIpAddress,omitempty"`
	UsePolicyBasedTrafficSelectors *bool                                       `json:"usePolicyBasedTrafficSelectors,omitempty"`
	VirtualNetworkGateway1         VirtualNetworkGateway                       `json:"virtualNetworkGateway1"`
	VirtualNetworkGateway2         *VirtualNetworkGateway                      `json:"virtualNetworkGateway2,omitempty"`
}
