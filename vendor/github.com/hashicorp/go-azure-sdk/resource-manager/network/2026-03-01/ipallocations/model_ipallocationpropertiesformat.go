package ipallocations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPAllocationPropertiesFormat struct {
	AllocationTags   *map[string]string `json:"allocationTags,omitempty"`
	IPamAllocationId *string            `json:"ipamAllocationId,omitempty"`
	Prefix           *string            `json:"prefix,omitempty"`
	PrefixLength     *int64             `json:"prefixLength,omitempty"`
	PrefixType       *IPVersion         `json:"prefixType,omitempty"`
	Subnet           *SubResource       `json:"subnet,omitempty"`
	Type             *IPAllocationType  `json:"type,omitempty"`
	VirtualNetwork   *SubResource       `json:"virtualNetwork,omitempty"`
}
