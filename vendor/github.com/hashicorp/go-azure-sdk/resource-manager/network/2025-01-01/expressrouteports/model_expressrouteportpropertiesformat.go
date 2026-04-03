package expressrouteports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRoutePortPropertiesFormat struct {
	AllocationDate             *string                         `json:"allocationDate,omitempty"`
	BandwidthInGbps            *int64                          `json:"bandwidthInGbps,omitempty"`
	BillingType                *ExpressRoutePortsBillingType   `json:"billingType,omitempty"`
	Circuits                   *[]SubResource                  `json:"circuits,omitempty"`
	Encapsulation              *ExpressRoutePortsEncapsulation `json:"encapsulation,omitempty"`
	EtherType                  *string                         `json:"etherType,omitempty"`
	Links                      *[]ExpressRouteLink             `json:"links,omitempty"`
	Mtu                        *string                         `json:"mtu,omitempty"`
	PeeringLocation            *string                         `json:"peeringLocation,omitempty"`
	ProvisionedBandwidthInGbps *float64                        `json:"provisionedBandwidthInGbps,omitempty"`
	ProvisioningState          *ProvisioningState              `json:"provisioningState,omitempty"`
	ResourceGuid               *string                         `json:"resourceGuid,omitempty"`
}
