package machineextensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineExtensionProperties struct {
	AutoUpgradeMinorVersion *bool                         `json:"autoUpgradeMinorVersion,omitempty"`
	EnableAutomaticUpgrade  *bool                         `json:"enableAutomaticUpgrade,omitempty"`
	ForceUpdateTag          *string                       `json:"forceUpdateTag,omitempty"`
	InstanceView            *MachineExtensionInstanceView `json:"instanceView,omitempty"`
	ProtectedSettings       *interface{}                  `json:"protectedSettings,omitempty"`
	ProvisioningState       *string                       `json:"provisioningState,omitempty"`
	Publisher               *string                       `json:"publisher,omitempty"`
	Settings                *interface{}                  `json:"settings,omitempty"`
	Type                    *string                       `json:"type,omitempty"`
	TypeHandlerVersion      *string                       `json:"typeHandlerVersion,omitempty"`
}
