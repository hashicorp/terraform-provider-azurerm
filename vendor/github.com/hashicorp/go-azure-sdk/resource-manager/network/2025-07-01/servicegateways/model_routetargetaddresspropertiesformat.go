package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteTargetAddressPropertiesFormat struct {
	PrivateIPAddress          *string             `json:"privateIPAddress,omitempty"`
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	Subnet                    *Subnet             `json:"subnet,omitempty"`
}
