package virtualnetworkappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkAppliancePropertiesFormat struct {
	BandwidthInGbps         *float64                                  `json:"bandwidthInGbps,omitempty"`
	IPConfigurations        *[]VirtualNetworkApplianceIPConfiguration `json:"ipConfigurations,omitempty"`
	PrivateIPAddressVersion *VirtualNetworkApplianceIPVersionType     `json:"privateIPAddressVersion,omitempty"`
	ProvisioningState       *ProvisioningState                        `json:"provisioningState,omitempty"`
	ResourceGuid            *string                                   `json:"resourceGuid,omitempty"`
	Subnet                  *Subnet                                   `json:"subnet,omitempty"`
}
