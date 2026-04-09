package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayFrontendIPConfigurationPropertiesFormat struct {
	PrivateIPAddress          *string             `json:"privateIPAddress,omitempty"`
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	PrivateLinkConfiguration  *SubResource        `json:"privateLinkConfiguration,omitempty"`
	ProvisioningState         *ProvisioningState  `json:"provisioningState,omitempty"`
	PublicIPAddress           *SubResource        `json:"publicIPAddress,omitempty"`
	Subnet                    *SubResource        `json:"subnet,omitempty"`
}
