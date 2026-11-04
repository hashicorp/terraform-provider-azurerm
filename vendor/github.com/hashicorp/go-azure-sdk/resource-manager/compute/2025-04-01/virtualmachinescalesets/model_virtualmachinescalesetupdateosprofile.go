package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetUpdateOSProfile struct {
	CustomData           *string               `json:"customData,omitempty"`
	LinuxConfiguration   *LinuxConfiguration   `json:"linuxConfiguration,omitempty"`
	Secrets              *[]VaultSecretGroup   `json:"secrets,omitempty"`
	WindowsConfiguration *WindowsConfiguration `json:"windowsConfiguration,omitempty"`
}
