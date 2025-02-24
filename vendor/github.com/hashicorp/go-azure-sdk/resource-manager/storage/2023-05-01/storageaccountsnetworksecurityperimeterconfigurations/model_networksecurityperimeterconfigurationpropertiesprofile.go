package storageaccountsnetworksecurityperimeterconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterConfigurationPropertiesProfile struct {
	AccessRules               *[]NspAccessRule `json:"accessRules,omitempty"`
	AccessRulesVersion        *float64         `json:"accessRulesVersion,omitempty"`
	DiagnosticSettingsVersion *float64         `json:"diagnosticSettingsVersion,omitempty"`
	EnabledLogCategories      *[]string        `json:"enabledLogCategories,omitempty"`
	Name                      *string          `json:"name,omitempty"`
}
