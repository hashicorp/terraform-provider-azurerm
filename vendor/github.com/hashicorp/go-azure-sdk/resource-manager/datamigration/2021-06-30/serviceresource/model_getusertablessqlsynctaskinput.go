package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUserTablesSqlSyncTaskInput struct {
	SelectedSourceDatabases []string          `json:"selectedSourceDatabases"`
	SelectedTargetDatabases []string          `json:"selectedTargetDatabases"`
	SourceConnectionInfo    SqlConnectionInfo `json:"sourceConnectionInfo"`
	TargetConnectionInfo    SqlConnectionInfo `json:"targetConnectionInfo"`
}
