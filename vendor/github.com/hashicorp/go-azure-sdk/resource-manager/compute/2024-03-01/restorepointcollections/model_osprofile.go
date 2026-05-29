package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSProfile struct {
	AdminPassword               *string               `json:"adminPassword,omitempty"`
	AdminUsername               *string               `json:"adminUsername,omitempty"`
	AllowExtensionOperations    *bool                 `json:"allowExtensionOperations,omitempty"`
	ComputerName                *string               `json:"computerName,omitempty"`
	CustomData                  *string               `json:"customData,omitempty"`
	LinuxConfiguration          *LinuxConfiguration   `json:"linuxConfiguration,omitempty"`
	RequireGuestProvisionSignal *bool                 `json:"requireGuestProvisionSignal,omitempty"`
	Secrets                     *[]VaultSecretGroup   `json:"secrets,omitempty"`
	WindowsConfiguration        *WindowsConfiguration `json:"windowsConfiguration,omitempty"`
}
