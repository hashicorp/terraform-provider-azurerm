package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceUpdate struct {
	IPv4AddressType  *AllocationMethod `json:"ipv4AddressType,omitempty"`
	IPv6AddressType  *AllocationMethod `json:"ipv6AddressType,omitempty"`
	MacAddress       *string           `json:"macAddress,omitempty"`
	MacAddressType   *AllocationMethod `json:"macAddressType,omitempty"`
	Name             *string           `json:"name,omitempty"`
	NicId            *string           `json:"nicId,omitempty"`
	VirtualNetworkId *string           `json:"virtualNetworkId,omitempty"`
}
