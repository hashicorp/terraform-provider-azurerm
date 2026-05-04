package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstallPatchesParameters struct {
	LinuxParameters   *LinuxParameters          `json:"linuxParameters,omitempty"`
	MaximumDuration   *string                   `json:"maximumDuration,omitempty"`
	RebootSetting     VMGuestPatchRebootSetting `json:"rebootSetting"`
	WindowsParameters *WindowsParameters        `json:"windowsParameters,omitempty"`
}
