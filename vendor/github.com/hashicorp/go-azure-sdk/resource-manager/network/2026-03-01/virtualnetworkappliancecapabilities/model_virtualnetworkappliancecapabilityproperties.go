package virtualnetworkappliancecapabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkApplianceCapabilityProperties struct {
	IPVersion         *VirtualNetworkApplianceCapabilityIPVersion `json:"ipVersion,omitempty"`
	LinkedResourceId  *string                                     `json:"linkedResourceId,omitempty"`
	ProvisioningState *ProvisioningState                          `json:"provisioningState,omitempty"`
}
