package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnoseRequestProperties struct {
	ApplicationInsights *map[string]interface{} `json:"applicationInsights,omitempty"`
	ContainerRegistry   *map[string]interface{} `json:"containerRegistry,omitempty"`
	DnsResolution       *map[string]interface{} `json:"dnsResolution,omitempty"`
	KeyVault            *map[string]interface{} `json:"keyVault,omitempty"`
	Nsg                 *map[string]interface{} `json:"nsg,omitempty"`
	Others              *map[string]interface{} `json:"others,omitempty"`
	ResourceLock        *map[string]interface{} `json:"resourceLock,omitempty"`
	StorageAccount      *map[string]interface{} `json:"storageAccount,omitempty"`
	Udr                 *map[string]interface{} `json:"udr,omitempty"`
}
