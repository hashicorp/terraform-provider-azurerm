package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlServerSqlMISyncTaskInput struct {
	AzureApp             AzureActiveDirectoryApp              `json:"azureApp"`
	BackupFileShare      *FileShare                           `json:"backupFileShare,omitempty"`
	SelectedDatabases    []MigrateSqlServerSqlMIDatabaseInput `json:"selectedDatabases"`
	SourceConnectionInfo SqlConnectionInfo                    `json:"sourceConnectionInfo"`
	StorageResourceId    string                               `json:"storageResourceId"`
	TargetConnectionInfo MiSqlConnectionInfo                  `json:"targetConnectionInfo"`
}
