package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstancePropertiesOsProfile struct {
	AdminPassword        *string                                                        `json:"adminPassword,omitempty"`
	AdminUsername        *string                                                        `json:"adminUsername,omitempty"`
	ComputerName         *string                                                        `json:"computerName,omitempty"`
	LinuxConfiguration   *VirtualMachineInstancePropertiesOsProfileLinuxConfiguration   `json:"linuxConfiguration,omitempty"`
	WindowsConfiguration *VirtualMachineInstancePropertiesOsProfileWindowsConfiguration `json:"windowsConfiguration,omitempty"`
}
