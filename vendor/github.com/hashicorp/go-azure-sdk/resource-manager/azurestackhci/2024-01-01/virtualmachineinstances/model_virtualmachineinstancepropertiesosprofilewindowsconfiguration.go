package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstancePropertiesOsProfileWindowsConfiguration struct {
	EnableAutomaticUpdates *bool             `json:"enableAutomaticUpdates,omitempty"`
	ProvisionVMAgent       *bool             `json:"provisionVMAgent,omitempty"`
	ProvisionVMConfigAgent *bool             `json:"provisionVMConfigAgent,omitempty"`
	Ssh                    *SshConfiguration `json:"ssh,omitempty"`
	TimeZone               *string           `json:"timeZone,omitempty"`
}
