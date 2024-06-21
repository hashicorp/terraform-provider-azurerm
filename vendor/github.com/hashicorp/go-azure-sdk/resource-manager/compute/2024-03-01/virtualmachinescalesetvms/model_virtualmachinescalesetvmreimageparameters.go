package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetVMReimageParameters struct {
	ExactVersion                  *string                    `json:"exactVersion,omitempty"`
	ForceUpdateOSDiskForEphemeral *bool                      `json:"forceUpdateOSDiskForEphemeral,omitempty"`
	OsProfile                     *OSProfileProvisioningData `json:"osProfile,omitempty"`
	TempDisk                      *bool                      `json:"tempDisk,omitempty"`
}
