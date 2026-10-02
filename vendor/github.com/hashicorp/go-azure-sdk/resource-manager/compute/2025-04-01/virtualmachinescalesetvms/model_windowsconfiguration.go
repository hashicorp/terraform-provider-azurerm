package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WindowsConfiguration struct {
	AdditionalUnattendContent    *[]AdditionalUnattendContent `json:"additionalUnattendContent,omitempty"`
	EnableAutomaticUpdates       *bool                        `json:"enableAutomaticUpdates,omitempty"`
	EnableVMAgentPlatformUpdates *bool                        `json:"enableVMAgentPlatformUpdates,omitempty"`
	PatchSettings                *PatchSettings               `json:"patchSettings,omitempty"`
	ProvisionVMAgent             *bool                        `json:"provisionVMAgent,omitempty"`
	TimeZone                     *string                      `json:"timeZone,omitempty"`
	WinRM                        *WinRMConfiguration          `json:"winRM,omitempty"`
}
