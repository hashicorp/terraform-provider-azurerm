package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSqlServerSqlMITaskInput struct {
	AadDomainName        *string                              `json:"aadDomainName,omitempty"`
	BackupBlobShare      BlobShare                            `json:"backupBlobShare"`
	BackupFileShare      *FileShare                           `json:"backupFileShare,omitempty"`
	BackupMode           *BackupMode                          `json:"backupMode,omitempty"`
	SelectedAgentJobs    *[]string                            `json:"selectedAgentJobs,omitempty"`
	SelectedDatabases    []MigrateSqlServerSqlMIDatabaseInput `json:"selectedDatabases"`
	SelectedLogins       *[]string                            `json:"selectedLogins,omitempty"`
	SourceConnectionInfo SqlConnectionInfo                    `json:"sourceConnectionInfo"`
	TargetConnectionInfo SqlConnectionInfo                    `json:"targetConnectionInfo"`
}
