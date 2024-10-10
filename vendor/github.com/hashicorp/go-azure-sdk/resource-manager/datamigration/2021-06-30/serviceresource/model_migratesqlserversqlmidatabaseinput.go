package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlMIDatabaseInput struct {
	BackupFilePaths     *[]string  `json:"backupFilePaths,omitempty"`
	BackupFileShare     *FileShare `json:"backupFileShare,omitempty"`
	Name                string     `json:"name"`
	RestoreDatabaseName string     `json:"restoreDatabaseName"`
}
