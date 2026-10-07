package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlAzureDbForMySqlOfflineDatabaseInput struct {
	Name               *string            `json:"name,omitempty"`
	TableMap           *map[string]string `json:"tableMap,omitempty"`
	TargetDatabaseName *string            `json:"targetDatabaseName,omitempty"`
}
