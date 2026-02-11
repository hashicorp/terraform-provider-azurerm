package networkprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPamPoolPrefixAllocation struct {
	AllocatedAddressPrefixes *[]string                     `json:"allocatedAddressPrefixes,omitempty"`
	NumberOfIPAddresses      *string                       `json:"numberOfIpAddresses,omitempty"`
	Pool                     *IPamPoolPrefixAllocationPool `json:"pool,omitempty"`
}
