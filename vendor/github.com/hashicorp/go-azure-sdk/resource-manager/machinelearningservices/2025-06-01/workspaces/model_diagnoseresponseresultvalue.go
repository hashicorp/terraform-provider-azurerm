package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnoseResponseResultValue struct {
	ApplicationInsightsResults *[]DiagnoseResult `json:"applicationInsightsResults,omitempty"`
	ContainerRegistryResults   *[]DiagnoseResult `json:"containerRegistryResults,omitempty"`
	DnsResolutionResults       *[]DiagnoseResult `json:"dnsResolutionResults,omitempty"`
	KeyVaultResults            *[]DiagnoseResult `json:"keyVaultResults,omitempty"`
	NetworkSecurityRuleResults *[]DiagnoseResult `json:"networkSecurityRuleResults,omitempty"`
	OtherResults               *[]DiagnoseResult `json:"otherResults,omitempty"`
	ResourceLockResults        *[]DiagnoseResult `json:"resourceLockResults,omitempty"`
	StorageAccountResults      *[]DiagnoseResult `json:"storageAccountResults,omitempty"`
	UserDefinedRouteResults    *[]DiagnoseResult `json:"userDefinedRouteResults,omitempty"`
}
