package routetables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceIPConfigurationProperties struct {
	Primary                   *bool               `json:"primary,omitempty"`
	PrivateIPAddress          *string             `json:"privateIPAddress,omitempty"`
	PrivateIPAddressVersion   *IPVersion          `json:"privateIPAddressVersion,omitempty"`
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         *ProvisioningState  `json:"provisioningState,omitempty"`
	Subnet                    *Subnet             `json:"subnet,omitempty"`
}
