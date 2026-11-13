package expressroutelags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagPropertiesFormat struct {
	AllocationDate             *string                         `json:"allocationDate,omitempty"`
	Authorizations             *[]ExpressRouteLagAuthorization `json:"authorizations,omitempty"`
	BandwidthInGbps            *int64                          `json:"bandwidthInGbps,omitempty"`
	BillingType                *ExpressRouteLagBillingType     `json:"billingType,omitempty"`
	Circuits                   *[]SubResource                  `json:"circuits,omitempty"`
	Encapsulation              *ExpressRouteLagEncapsulation   `json:"encapsulation,omitempty"`
	EtherType                  *string                         `json:"etherType,omitempty"`
	LacpTimer                  *ExpressRouteLagLacpTimer       `json:"lacpTimer,omitempty"`
	Links                      *[]ExpressRouteLagLink          `json:"links,omitempty"`
	MinimumActivePortsRequired *int64                          `json:"minimumActivePortsRequired,omitempty"`
	Mtu                        *string                         `json:"mtu,omitempty"`
	NumberOfPorts              *int64                          `json:"numberOfPorts,omitempty"`
	PeeringLocation            *string                         `json:"peeringLocation,omitempty"`
	ProvisionedBandwidthInGbps *float64                        `json:"provisionedBandwidthInGbps,omitempty"`
	ProvisioningState          *ProvisioningState              `json:"provisioningState,omitempty"`
	ResourceGuid               *string                         `json:"resourceGuid,omitempty"`
}
