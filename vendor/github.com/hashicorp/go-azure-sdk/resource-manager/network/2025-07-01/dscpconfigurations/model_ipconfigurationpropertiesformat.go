package dscpconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigurationPropertiesFormat struct {
	PrivateIPAddress          *string             `json:"privateIPAddress,omitempty"`
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         *ProvisioningState  `json:"provisioningState,omitempty"`
	PublicIPAddress           *PublicIPAddress    `json:"publicIPAddress,omitempty"`
	Subnet                    *Subnet             `json:"subnet,omitempty"`
}
