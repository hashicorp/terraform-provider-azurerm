package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolUpdateProperties struct {
	DevBoxDefinition             *PoolDevBoxDefinition          `json:"devBoxDefinition,omitempty"`
	DevBoxDefinitionName         *string                        `json:"devBoxDefinitionName,omitempty"`
	DevBoxDefinitionType         *PoolDevBoxDefinitionType      `json:"devBoxDefinitionType,omitempty"`
	DisplayName                  *string                        `json:"displayName,omitempty"`
	LicenseType                  *LicenseType                   `json:"licenseType,omitempty"`
	LocalAdministrator           *LocalAdminStatus              `json:"localAdministrator,omitempty"`
	ManagedVirtualNetworkRegions *[]string                      `json:"managedVirtualNetworkRegions,omitempty"`
	NetworkConnectionName        *string                        `json:"networkConnectionName,omitempty"`
	SingleSignOnStatus           *SingleSignOnStatus            `json:"singleSignOnStatus,omitempty"`
	StopOnDisconnect             *StopOnDisconnectConfiguration `json:"stopOnDisconnect,omitempty"`
	StopOnNoConnect              *StopOnNoConnectConfiguration  `json:"stopOnNoConnect,omitempty"`
	VirtualNetworkType           *VirtualNetworkType            `json:"virtualNetworkType,omitempty"`
}
