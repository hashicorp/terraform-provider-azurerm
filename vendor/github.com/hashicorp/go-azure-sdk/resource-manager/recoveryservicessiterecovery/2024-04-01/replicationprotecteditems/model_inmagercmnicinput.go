package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmNicInput struct {
	IsPrimaryNic          string  `json:"isPrimaryNic"`
	IsSelectedForFailover *string `json:"isSelectedForFailover,omitempty"`
	NicId                 string  `json:"nicId"`
	TargetStaticIPAddress *string `json:"targetStaticIPAddress,omitempty"`
	TargetSubnetName      *string `json:"targetSubnetName,omitempty"`
	TestStaticIPAddress   *string `json:"testStaticIPAddress,omitempty"`
	TestSubnetName        *string `json:"testSubnetName,omitempty"`
}
