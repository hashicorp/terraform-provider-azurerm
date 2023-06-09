package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionProperties struct {
	BandwidthInMbps                  *int64                                `json:"bandwidthInMbps,omitempty"`
	ExpressRouteCircuit              *ExpressRouteCircuitReference         `json:"expressRouteCircuit,omitempty"`
	PeeringLocation                  *string                               `json:"peeringLocation,omitempty"`
	Peerings                         *[]ExpressRouteCrossConnectionPeering `json:"peerings,omitempty"`
	PrimaryAzurePort                 *string                               `json:"primaryAzurePort,omitempty"`
	ProvisioningState                *ProvisioningState                    `json:"provisioningState,omitempty"`
	STag                             *int64                                `json:"sTag,omitempty"`
	SecondaryAzurePort               *string                               `json:"secondaryAzurePort,omitempty"`
	ServiceProviderNotes             *string                               `json:"serviceProviderNotes,omitempty"`
	ServiceProviderProvisioningState *ServiceProviderProvisioningState     `json:"serviceProviderProvisioningState,omitempty"`
}
