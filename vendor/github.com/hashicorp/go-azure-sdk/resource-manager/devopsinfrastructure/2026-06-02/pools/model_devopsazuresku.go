package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevOpsAzureSku struct {
	LinuxNVMePath    *string   `json:"linuxNvmePath,omitempty"`
	Name             string    `json:"name"`
	VMSizes          *[]VMSize `json:"vmSizes,omitempty"`
	WindowsNVMeDrive *string   `json:"windowsNvmeDrive,omitempty"`
}
