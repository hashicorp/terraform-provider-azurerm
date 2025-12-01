package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPropertiesInfraVnetProfileHci struct {
	MocGroup    *string `json:"mocGroup,omitempty"`
	MocLocation *string `json:"mocLocation,omitempty"`
	MocVnetName *string `json:"mocVnetName,omitempty"`
}
