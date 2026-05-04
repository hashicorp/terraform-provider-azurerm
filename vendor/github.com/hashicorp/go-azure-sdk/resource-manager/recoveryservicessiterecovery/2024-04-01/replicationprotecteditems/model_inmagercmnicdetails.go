package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmNicDetails struct {
	IsPrimaryNic          *string              `json:"isPrimaryNic,omitempty"`
	IsSelectedForFailover *string              `json:"isSelectedForFailover,omitempty"`
	NicId                 *string              `json:"nicId,omitempty"`
	SourceIPAddress       *string              `json:"sourceIPAddress,omitempty"`
	SourceIPAddressType   *EthernetAddressType `json:"sourceIPAddressType,omitempty"`
	SourceNetworkId       *string              `json:"sourceNetworkId,omitempty"`
	SourceSubnetName      *string              `json:"sourceSubnetName,omitempty"`
	TargetIPAddress       *string              `json:"targetIPAddress,omitempty"`
	TargetIPAddressType   *EthernetAddressType `json:"targetIPAddressType,omitempty"`
	TargetSubnetName      *string              `json:"targetSubnetName,omitempty"`
	TestIPAddress         *string              `json:"testIPAddress,omitempty"`
	TestIPAddressType     *EthernetAddressType `json:"testIPAddressType,omitempty"`
	TestSubnetName        *string              `json:"testSubnetName,omitempty"`
}
