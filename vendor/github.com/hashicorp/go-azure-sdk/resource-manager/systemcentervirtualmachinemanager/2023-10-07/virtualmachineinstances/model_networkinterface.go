package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterface struct {
	DisplayName      *string           `json:"displayName,omitempty"`
	IPv4AddressType  *AllocationMethod `json:"ipv4AddressType,omitempty"`
	IPv4Addresses    *[]string         `json:"ipv4Addresses,omitempty"`
	IPv6AddressType  *AllocationMethod `json:"ipv6AddressType,omitempty"`
	IPv6Addresses    *[]string         `json:"ipv6Addresses,omitempty"`
	MacAddress       *string           `json:"macAddress,omitempty"`
	MacAddressType   *AllocationMethod `json:"macAddressType,omitempty"`
	Name             *string           `json:"name,omitempty"`
	NetworkName      *string           `json:"networkName,omitempty"`
	NicId            *string           `json:"nicId,omitempty"`
	VirtualNetworkId *string           `json:"virtualNetworkId,omitempty"`
}
