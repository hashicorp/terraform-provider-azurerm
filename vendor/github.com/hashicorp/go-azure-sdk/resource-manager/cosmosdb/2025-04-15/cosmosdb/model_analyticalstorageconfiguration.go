<<<<<<<< HEAD:vendor/github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems/model_amlfilesystemarchiveinfo.go
package amlfilesystems
========
package cosmosdb
>>>>>>>> 3e1ec57095 (cosmosdb: upgrade API to 2025-04-15):vendor/github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-04-15/cosmosdb/model_analyticalstorageconfiguration.go

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

<<<<<<<< HEAD:vendor/github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems/model_amlfilesystemarchiveinfo.go
type AmlFilesystemArchiveInfo struct {
	FilesystemPath *string `json:"filesystemPath,omitempty"`
========
type AnalyticalStorageConfiguration struct {
	SchemaType *AnalyticalStorageSchemaType `json:"schemaType,omitempty"`
>>>>>>>> 3e1ec57095 (cosmosdb: upgrade API to 2025-04-15):vendor/github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-04-15/cosmosdb/model_analyticalstorageconfiguration.go
}
