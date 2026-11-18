package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreRequestProperties struct {
	AdjustConnectionStrings    *bool                       `json:"adjustConnectionStrings,omitempty"`
	AppServicePlan             *string                     `json:"appServicePlan,omitempty"`
	BlobName                   *string                     `json:"blobName,omitempty"`
	Databases                  *[]DatabaseBackupSetting    `json:"databases,omitempty"`
	HostingEnvironment         *string                     `json:"hostingEnvironment,omitempty"`
	IgnoreConflictingHostNames *bool                       `json:"ignoreConflictingHostNames,omitempty"`
	IgnoreDatabases            *bool                       `json:"ignoreDatabases,omitempty"`
	OperationType              *BackupRestoreOperationType `json:"operationType,omitempty"`
	Overwrite                  bool                        `json:"overwrite"`
	SiteName                   *string                     `json:"siteName,omitempty"`
	StorageAccountURL          string                      `json:"storageAccountUrl"`
}
