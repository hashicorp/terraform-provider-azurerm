package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineReimageParameters struct {
	ExactVersion *string                    `json:"exactVersion,omitempty"`
	OsProfile    *OSProfileProvisioningData `json:"osProfile,omitempty"`
	TempDisk     *bool                      `json:"tempDisk,omitempty"`
}
