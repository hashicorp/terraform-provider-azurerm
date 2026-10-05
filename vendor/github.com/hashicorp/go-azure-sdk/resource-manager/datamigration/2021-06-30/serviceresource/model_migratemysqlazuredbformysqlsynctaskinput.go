package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlAzureDbForMySqlSyncTaskInput struct {
	SelectedDatabases    []MigrateMySqlAzureDbForMySqlSyncDatabaseInput `json:"selectedDatabases"`
	SourceConnectionInfo MySqlConnectionInfo                            `json:"sourceConnectionInfo"`
	TargetConnectionInfo MySqlConnectionInfo                            `json:"targetConnectionInfo"`
}
