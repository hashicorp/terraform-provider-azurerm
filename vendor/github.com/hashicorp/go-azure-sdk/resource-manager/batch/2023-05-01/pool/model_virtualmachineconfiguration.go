package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineConfiguration struct {
	ContainerConfiguration      *ContainerConfiguration      `json:"containerConfiguration,omitempty"`
	DataDisks                   *[]DataDisk                  `json:"dataDisks,omitempty"`
	DiskEncryptionConfiguration *DiskEncryptionConfiguration `json:"diskEncryptionConfiguration,omitempty"`
	Extensions                  *[]VmExtension               `json:"extensions,omitempty"`
	ImageReference              ImageReference               `json:"imageReference"`
	LicenseType                 *string                      `json:"licenseType,omitempty"`
	NodeAgentSkuId              string                       `json:"nodeAgentSkuId"`
	NodePlacementConfiguration  *NodePlacementConfiguration  `json:"nodePlacementConfiguration,omitempty"`
	OsDisk                      *OSDisk                      `json:"osDisk,omitempty"`
	WindowsConfiguration        *WindowsConfiguration        `json:"windowsConfiguration,omitempty"`
}
