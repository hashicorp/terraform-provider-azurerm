package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorePointSourceVMStorageProfile struct {
	DataDisks          *[]RestorePointSourceVMDataDisk `json:"dataDisks,omitempty"`
	DiskControllerType *DiskControllerTypes            `json:"diskControllerType,omitempty"`
	OsDisk             *RestorePointSourceVMOSDisk     `json:"osDisk,omitempty"`
}
