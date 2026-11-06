package virtualmachineimagetemplate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateVMProfile struct {
	OsDiskSizeGB           *int64                `json:"osDiskSizeGB,omitempty"`
	UserAssignedIdentities *[]string             `json:"userAssignedIdentities,omitempty"`
	VMSize                 *string               `json:"vmSize,omitempty"`
	VnetConfig             *VirtualNetworkConfig `json:"vnetConfig,omitempty"`
}
