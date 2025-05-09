package localnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddressSpace struct {
	AddressPrefixes           *[]string                   `json:"addressPrefixes,omitempty"`
	IPamPoolPrefixAllocations *[]IPamPoolPrefixAllocation `json:"ipamPoolPrefixAllocations,omitempty"`
}
