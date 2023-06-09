package expressroutecircuitpeerings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeerExpressRouteCircuitConnectionPropertiesFormat struct {
	AddressPrefix                  *string                  `json:"addressPrefix,omitempty"`
	AuthResourceGuid               *string                  `json:"authResourceGuid,omitempty"`
	CircuitConnectionStatus        *CircuitConnectionStatus `json:"circuitConnectionStatus,omitempty"`
	ConnectionName                 *string                  `json:"connectionName,omitempty"`
	ExpressRouteCircuitPeering     *SubResource             `json:"expressRouteCircuitPeering,omitempty"`
	PeerExpressRouteCircuitPeering *SubResource             `json:"peerExpressRouteCircuitPeering,omitempty"`
	ProvisioningState              *ProvisioningState       `json:"provisioningState,omitempty"`
}
