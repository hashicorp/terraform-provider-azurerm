package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigratePostgreSqlAzureDbForPostgreSqlSyncTaskInput struct {
	SelectedDatabases    []MigratePostgreSqlAzureDbForPostgreSqlSyncDatabaseInput `json:"selectedDatabases"`
	SourceConnectionInfo PostgreSqlConnectionInfo                                 `json:"sourceConnectionInfo"`
	TargetConnectionInfo PostgreSqlConnectionInfo                                 `json:"targetConnectionInfo"`
}
