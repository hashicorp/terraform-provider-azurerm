package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenStore struct {
	AzureBlobStorage           *BlobStorageTokenStore `json:"azureBlobStorage,omitempty"`
	Enabled                    *bool                  `json:"enabled,omitempty"`
	FileSystem                 *FileSystemTokenStore  `json:"fileSystem,omitempty"`
	TokenRefreshExtensionHours *float64               `json:"tokenRefreshExtensionHours,omitempty"`
}
