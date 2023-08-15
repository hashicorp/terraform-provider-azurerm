package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionHostIPConfigurationPropertiesFormat struct {
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         *ProvisioningState  `json:"provisioningState,omitempty"`
	PublicIPAddress           SubResource         `json:"publicIPAddress"`
	Subnet                    SubResource         `json:"subnet"`
}
