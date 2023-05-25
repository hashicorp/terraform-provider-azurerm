package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccount struct {
	Container     *string `json:"container,omitempty"`
	FileSystem    *string `json:"fileSystem,omitempty"`
	Fileshare     *string `json:"fileshare,omitempty"`
	IsDefault     *bool   `json:"isDefault,omitempty"`
	Key           *string `json:"key,omitempty"`
	MsiResourceId *string `json:"msiResourceId,omitempty"`
	Name          *string `json:"name,omitempty"`
	ResourceId    *string `json:"resourceId,omitempty"`
	Saskey        *string `json:"saskey,omitempty"`
}
