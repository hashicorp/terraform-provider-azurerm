package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseConfiguration struct {
	DatabaseType                *SAPDatabaseType            `json:"databaseType,omitempty"`
	DiskConfiguration           *DiskConfiguration          `json:"diskConfiguration,omitempty"`
	InstanceCount               int64                       `json:"instanceCount"`
	SubnetId                    string                      `json:"subnetId"`
	VirtualMachineConfiguration VirtualMachineConfiguration `json:"virtualMachineConfiguration"`
}
