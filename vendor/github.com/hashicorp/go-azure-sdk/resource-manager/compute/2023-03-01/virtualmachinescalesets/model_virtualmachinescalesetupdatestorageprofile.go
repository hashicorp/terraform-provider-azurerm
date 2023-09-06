package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetUpdateStorageProfile struct {
	DataDisks          *[]VirtualMachineScaleSetDataDisk   `json:"dataDisks,omitempty"`
	DiskControllerType *string                             `json:"diskControllerType,omitempty"`
	ImageReference     *ImageReference                     `json:"imageReference,omitempty"`
	OsDisk             *VirtualMachineScaleSetUpdateOSDisk `json:"osDisk,omitempty"`
}
