package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigInputDetails struct {
	IPConfigName                    *string   `json:"ipConfigName,omitempty"`
	IsPrimary                       *bool     `json:"isPrimary,omitempty"`
	IsSeletedForFailover            *bool     `json:"isSeletedForFailover,omitempty"`
	RecoveryLBBackendAddressPoolIds *[]string `json:"recoveryLBBackendAddressPoolIds,omitempty"`
	RecoveryPublicIPAddressId       *string   `json:"recoveryPublicIPAddressId,omitempty"`
	RecoveryStaticIPAddress         *string   `json:"recoveryStaticIPAddress,omitempty"`
	RecoverySubnetName              *string   `json:"recoverySubnetName,omitempty"`
	TfoLBBackendAddressPoolIds      *[]string `json:"tfoLBBackendAddressPoolIds,omitempty"`
	TfoPublicIPAddressId            *string   `json:"tfoPublicIPAddressId,omitempty"`
	TfoStaticIPAddress              *string   `json:"tfoStaticIPAddress,omitempty"`
	TfoSubnetName                   *string   `json:"tfoSubnetName,omitempty"`
}
