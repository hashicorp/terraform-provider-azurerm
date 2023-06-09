package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitPropertiesFormat struct {
	AllowClassicOperations           *bool                                         `json:"allowClassicOperations,omitempty"`
	AuthorizationKey                 *string                                       `json:"authorizationKey,omitempty"`
	AuthorizationStatus              *string                                       `json:"authorizationStatus,omitempty"`
	Authorizations                   *[]ExpressRouteCircuitAuthorization           `json:"authorizations,omitempty"`
	BandwidthInGbps                  *float64                                      `json:"bandwidthInGbps,omitempty"`
	CircuitProvisioningState         *string                                       `json:"circuitProvisioningState,omitempty"`
	ExpressRoutePort                 *SubResource                                  `json:"expressRoutePort,omitempty"`
	GatewayManagerEtag               *string                                       `json:"gatewayManagerEtag,omitempty"`
	GlobalReachEnabled               *bool                                         `json:"globalReachEnabled,omitempty"`
	Peerings                         *[]ExpressRouteCircuitPeering                 `json:"peerings,omitempty"`
	ProvisioningState                *ProvisioningState                            `json:"provisioningState,omitempty"`
	ServiceKey                       *string                                       `json:"serviceKey,omitempty"`
	ServiceProviderNotes             *string                                       `json:"serviceProviderNotes,omitempty"`
	ServiceProviderProperties        *ExpressRouteCircuitServiceProviderProperties `json:"serviceProviderProperties,omitempty"`
	ServiceProviderProvisioningState *ServiceProviderProvisioningState             `json:"serviceProviderProvisioningState,omitempty"`
	Stag                             *int64                                        `json:"stag,omitempty"`
}
