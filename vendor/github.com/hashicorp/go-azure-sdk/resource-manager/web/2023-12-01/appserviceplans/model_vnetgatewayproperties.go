package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetGatewayProperties struct {
	VnetName      *string `json:"vnetName,omitempty"`
	VpnPackageUri string  `json:"vpnPackageUri"`
}
