package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitPeeringPropertiesFormat struct {
	AzureASN                   *int64                                `json:"azureASN,omitempty"`
	Connections                *[]ExpressRouteCircuitConnection      `json:"connections,omitempty"`
	ExpressRouteConnection     *ExpressRouteConnectionId             `json:"expressRouteConnection,omitempty"`
	GatewayManagerEtag         *string                               `json:"gatewayManagerEtag,omitempty"`
	IPv6PeeringConfig          *IPv6ExpressRouteCircuitPeeringConfig `json:"ipv6PeeringConfig,omitempty"`
	LastModifiedBy             *string                               `json:"lastModifiedBy,omitempty"`
	MicrosoftPeeringConfig     *ExpressRouteCircuitPeeringConfig     `json:"microsoftPeeringConfig,omitempty"`
	PeerASN                    *int64                                `json:"peerASN,omitempty"`
	PeeredConnections          *[]PeerExpressRouteCircuitConnection  `json:"peeredConnections,omitempty"`
	PeeringType                *ExpressRoutePeeringType              `json:"peeringType,omitempty"`
	PrimaryAzurePort           *string                               `json:"primaryAzurePort,omitempty"`
	PrimaryPeerAddressPrefix   *string                               `json:"primaryPeerAddressPrefix,omitempty"`
	ProvisioningState          *ProvisioningState                    `json:"provisioningState,omitempty"`
	RouteFilter                *SubResource                          `json:"routeFilter,omitempty"`
	SecondaryAzurePort         *string                               `json:"secondaryAzurePort,omitempty"`
	SecondaryPeerAddressPrefix *string                               `json:"secondaryPeerAddressPrefix,omitempty"`
	SharedKey                  *string                               `json:"sharedKey,omitempty"`
	State                      *ExpressRoutePeeringState             `json:"state,omitempty"`
	Stats                      *ExpressRouteCircuitStats             `json:"stats,omitempty"`
	VlanId                     *int64                                `json:"vlanId,omitempty"`
}
