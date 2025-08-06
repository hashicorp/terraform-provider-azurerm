package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceFile struct {
	AutoStorageContainerName *string                       `json:"autoStorageContainerName,omitempty"`
	BlobPrefix               *string                       `json:"blobPrefix,omitempty"`
	FileMode                 *string                       `json:"fileMode,omitempty"`
	FilePath                 *string                       `json:"filePath,omitempty"`
	HTTPURL                  *string                       `json:"httpUrl,omitempty"`
	IdentityReference        *ComputeNodeIdentityReference `json:"identityReference,omitempty"`
	StorageContainerURL      *string                       `json:"storageContainerUrl,omitempty"`
}
