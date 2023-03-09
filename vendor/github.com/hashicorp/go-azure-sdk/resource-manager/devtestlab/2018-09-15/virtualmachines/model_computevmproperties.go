package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeVMProperties struct {
	DataDiskIds        *[]string                      `json:"dataDiskIds,omitempty"`
	DataDisks          *[]ComputeDataDisk             `json:"dataDisks,omitempty"`
	NetworkInterfaceId *string                        `json:"networkInterfaceId,omitempty"`
	OsDiskId           *string                        `json:"osDiskId,omitempty"`
	OsType             *string                        `json:"osType,omitempty"`
	Statuses           *[]ComputeVMInstanceViewStatus `json:"statuses,omitempty"`
	VMSize             *string                        `json:"vmSize,omitempty"`
}
