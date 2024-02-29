package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetParametersProperties struct {
	SubnetResourceId  *string `json:"subnetResourceId,omitempty"`
	VnetName          *string `json:"vnetName,omitempty"`
	VnetResourceGroup *string `json:"vnetResourceGroup,omitempty"`
	VnetSubnetName    *string `json:"vnetSubnetName,omitempty"`
}
