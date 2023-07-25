package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigDetails struct {
	IPAddressType                   *string   `json:"ipAddressType,omitempty"`
	IsPrimary                       *bool     `json:"isPrimary,omitempty"`
	IsSeletedForFailover            *bool     `json:"isSeletedForFailover,omitempty"`
	Name                            *string   `json:"name,omitempty"`
	RecoveryIPAddressType           *string   `json:"recoveryIPAddressType,omitempty"`
	RecoveryLBBackendAddressPoolIds *[]string `json:"recoveryLBBackendAddressPoolIds,omitempty"`
	RecoveryPublicIPAddressId       *string   `json:"recoveryPublicIPAddressId,omitempty"`
	RecoveryStaticIPAddress         *string   `json:"recoveryStaticIPAddress,omitempty"`
	RecoverySubnetName              *string   `json:"recoverySubnetName,omitempty"`
	StaticIPAddress                 *string   `json:"staticIPAddress,omitempty"`
	SubnetName                      *string   `json:"subnetName,omitempty"`
	TfoLBBackendAddressPoolIds      *[]string `json:"tfoLBBackendAddressPoolIds,omitempty"`
	TfoPublicIPAddressId            *string   `json:"tfoPublicIPAddressId,omitempty"`
	TfoStaticIPAddress              *string   `json:"tfoStaticIPAddress,omitempty"`
	TfoSubnetName                   *string   `json:"tfoSubnetName,omitempty"`
}
