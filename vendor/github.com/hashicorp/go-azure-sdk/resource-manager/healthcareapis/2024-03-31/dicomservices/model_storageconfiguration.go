package dicomservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConfiguration struct {
	FileSystemName    *string `json:"fileSystemName,omitempty"`
	StorageResourceId *string `json:"storageResourceId,omitempty"`
}
