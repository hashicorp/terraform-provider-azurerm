package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigurationProperties struct {
	Primary                   *bool                                      `json:"primary,omitempty"`
	PrivateIPAddress          *string                                    `json:"privateIPAddress,omitempty"`
	PrivateIPAllocationMethod *PrivateIPAllocationMethod                 `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         *PrivateLinkConfigurationProvisioningState `json:"provisioningState,omitempty"`
	Subnet                    *ResourceId                                `json:"subnet,omitempty"`
}
