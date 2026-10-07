package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VmExtension struct {
	AutoUpgradeMinorVersion  *bool        `json:"autoUpgradeMinorVersion,omitempty"`
	EnableAutomaticUpgrade   *bool        `json:"enableAutomaticUpgrade,omitempty"`
	Name                     string       `json:"name"`
	ProtectedSettings        *interface{} `json:"protectedSettings,omitempty"`
	ProvisionAfterExtensions *[]string    `json:"provisionAfterExtensions,omitempty"`
	Publisher                string       `json:"publisher"`
	Settings                 *interface{} `json:"settings,omitempty"`
	Type                     string       `json:"type"`
	TypeHandlerVersion       *string      `json:"typeHandlerVersion,omitempty"`
}
