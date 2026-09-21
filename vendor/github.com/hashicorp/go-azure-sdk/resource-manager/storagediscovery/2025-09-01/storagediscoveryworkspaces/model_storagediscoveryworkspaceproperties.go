package storagediscoveryworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDiscoveryWorkspaceProperties struct {
	Description       *string                    `json:"description,omitempty"`
	ProvisioningState *ResourceProvisioningState `json:"provisioningState,omitempty"`
	Scopes            []StorageDiscoveryScope    `json:"scopes"`
	Sku               *StorageDiscoverySku       `json:"sku,omitempty"`
	WorkspaceRoots    []string                   `json:"workspaceRoots"`
}
