package routefilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitConnectionPropertiesFormat struct {
	AddressPrefix                  *string                      `json:"addressPrefix,omitempty"`
	AuthorizationKey               *string                      `json:"authorizationKey,omitempty"`
	CircuitConnectionStatus        *CircuitConnectionStatus     `json:"circuitConnectionStatus,omitempty"`
	ExpressRouteCircuitPeering     *SubResource                 `json:"expressRouteCircuitPeering,omitempty"`
	IPv6CircuitConnectionConfig    *IPv6CircuitConnectionConfig `json:"ipv6CircuitConnectionConfig,omitempty"`
	PeerExpressRouteCircuitPeering *SubResource                 `json:"peerExpressRouteCircuitPeering,omitempty"`
	ProvisioningState              *ProvisioningState           `json:"provisioningState,omitempty"`
}
