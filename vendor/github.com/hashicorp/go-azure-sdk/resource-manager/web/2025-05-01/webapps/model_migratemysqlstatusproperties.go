package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateMySqlStatusProperties struct {
	LocalMySqlEnabled        *bool            `json:"localMySqlEnabled,omitempty"`
	MigrationOperationStatus *OperationStatus `json:"migrationOperationStatus,omitempty"`
	OperationId              *string          `json:"operationId,omitempty"`
}
