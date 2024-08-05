package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationLogsConfig struct {
	AzureBlobStorage  *AzureBlobStorageApplicationLogsConfig  `json:"azureBlobStorage,omitempty"`
	AzureTableStorage *AzureTableStorageApplicationLogsConfig `json:"azureTableStorage,omitempty"`
	FileSystem        *FileSystemApplicationLogsConfig        `json:"fileSystem,omitempty"`
}
