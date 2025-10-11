package ipampools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolUsage struct {
	AddressPrefixes              *[]string         `json:"addressPrefixes,omitempty"`
	AllocatedAddressPrefixes     *[]string         `json:"allocatedAddressPrefixes,omitempty"`
	AvailableAddressPrefixes     *[]string         `json:"availableAddressPrefixes,omitempty"`
	ChildPools                   *[]ResourceBasics `json:"childPools,omitempty"`
	NumberOfAllocatedIPAddresses *string           `json:"numberOfAllocatedIPAddresses,omitempty"`
	NumberOfAvailableIPAddresses *string           `json:"numberOfAvailableIPAddresses,omitempty"`
	NumberOfReservedIPAddresses  *string           `json:"numberOfReservedIPAddresses,omitempty"`
	ReservedAddressPrefixes      *[]string         `json:"reservedAddressPrefixes,omitempty"`
	TotalNumberOfIPAddresses     *string           `json:"totalNumberOfIPAddresses,omitempty"`
}
