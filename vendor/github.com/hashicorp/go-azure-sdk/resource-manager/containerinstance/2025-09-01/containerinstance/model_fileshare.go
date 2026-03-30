package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShare struct {
	Name               *string              `json:"name,omitempty"`
	Properties         *FileShareProperties `json:"properties,omitempty"`
	ResourceGroupName  *string              `json:"resourceGroupName,omitempty"`
	StorageAccountName *string              `json:"storageAccountName,omitempty"`
}
