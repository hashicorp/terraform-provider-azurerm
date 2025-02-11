package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateMigrationInputSqlServerSqlMITaskOutput struct {
	BackupFolderErrors           *[]ReportableException `json:"backupFolderErrors,omitempty"`
	BackupShareCredentialsErrors *[]ReportableException `json:"backupShareCredentialsErrors,omitempty"`
	BackupStorageAccountErrors   *[]ReportableException `json:"backupStorageAccountErrors,omitempty"`
	DatabaseBackupInfo           *DatabaseBackupInfo    `json:"databaseBackupInfo,omitempty"`
	ExistingBackupErrors         *[]ReportableException `json:"existingBackupErrors,omitempty"`
	Id                           *string                `json:"id,omitempty"`
	Name                         *string                `json:"name,omitempty"`
	RestoreDatabaseNameErrors    *[]ReportableException `json:"restoreDatabaseNameErrors,omitempty"`
}
