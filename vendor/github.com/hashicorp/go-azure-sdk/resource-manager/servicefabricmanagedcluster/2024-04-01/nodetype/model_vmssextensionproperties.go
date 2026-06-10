package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSSExtensionProperties struct {
	AutoUpgradeMinorVersion  *bool                      `json:"autoUpgradeMinorVersion,omitempty"`
	EnableAutomaticUpgrade   *bool                      `json:"enableAutomaticUpgrade,omitempty"`
	ForceUpdateTag           *string                    `json:"forceUpdateTag,omitempty"`
	ProtectedSettings        *interface{}               `json:"protectedSettings,omitempty"`
	ProvisionAfterExtensions *[]string                  `json:"provisionAfterExtensions,omitempty"`
	ProvisioningState        *string                    `json:"provisioningState,omitempty"`
	Publisher                string                     `json:"publisher"`
	Settings                 *interface{}               `json:"settings,omitempty"`
	SetupOrder               *[]VMSSExtensionSetupOrder `json:"setupOrder,omitempty"`
	Type                     string                     `json:"type"`
	TypeHandlerVersion       string                     `json:"typeHandlerVersion"`
}
