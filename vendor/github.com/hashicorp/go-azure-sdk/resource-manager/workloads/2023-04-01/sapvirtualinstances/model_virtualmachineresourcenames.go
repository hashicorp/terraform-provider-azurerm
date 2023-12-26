package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineResourceNames struct {
	DataDiskNames      *map[string][]string             `json:"dataDiskNames,omitempty"`
	HostName           *string                          `json:"hostName,omitempty"`
	NetworkInterfaces  *[]NetworkInterfaceResourceNames `json:"networkInterfaces,omitempty"`
	OsDiskName         *string                          `json:"osDiskName,omitempty"`
	VirtualMachineName *string                          `json:"vmName,omitempty"`
}
