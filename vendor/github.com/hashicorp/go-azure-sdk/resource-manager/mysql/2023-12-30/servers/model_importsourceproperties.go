package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportSourceProperties struct {
	DataDirPath *string                  `json:"dataDirPath,omitempty"`
	SasToken    *string                  `json:"sasToken,omitempty"`
	StorageType *ImportSourceStorageType `json:"storageType,omitempty"`
	StorageURL  *string                  `json:"storageUrl,omitempty"`
}
