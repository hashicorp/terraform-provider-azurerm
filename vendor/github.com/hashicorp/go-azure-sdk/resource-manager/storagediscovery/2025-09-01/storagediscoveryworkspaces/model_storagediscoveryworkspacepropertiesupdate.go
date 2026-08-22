package storagediscoveryworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDiscoveryWorkspacePropertiesUpdate struct {
	Description    *string                  `json:"description,omitempty"`
	Scopes         *[]StorageDiscoveryScope `json:"scopes,omitempty"`
	Sku            *StorageDiscoverySku     `json:"sku,omitempty"`
	WorkspaceRoots *[]string                `json:"workspaceRoots,omitempty"`
}
