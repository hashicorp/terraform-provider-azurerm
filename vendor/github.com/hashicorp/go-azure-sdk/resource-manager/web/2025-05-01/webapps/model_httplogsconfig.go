package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPLogsConfig struct {
	AzureBlobStorage *AzureBlobStorageHTTPLogsConfig `json:"azureBlobStorage,omitempty"`
	FileSystem       *FileSystemHTTPLogsConfig       `json:"fileSystem,omitempty"`
}
