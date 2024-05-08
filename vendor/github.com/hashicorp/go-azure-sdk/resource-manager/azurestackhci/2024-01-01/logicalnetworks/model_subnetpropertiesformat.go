package logicalnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetPropertiesFormat struct {
	AddressPrefix             *string                                                   `json:"addressPrefix,omitempty"`
	AddressPrefixes           *[]string                                                 `json:"addressPrefixes,omitempty"`
	IPAllocationMethod        *IPAllocationMethodEnum                                   `json:"ipAllocationMethod,omitempty"`
	IPConfigurationReferences *[]SubnetPropertiesFormatIPConfigurationReferencesInlined `json:"ipConfigurationReferences,omitempty"`
	IPPools                   *[]IPPool                                                 `json:"ipPools,omitempty"`
	RouteTable                *RouteTable                                               `json:"routeTable,omitempty"`
	Vlan                      *int64                                                    `json:"vlan,omitempty"`
}
