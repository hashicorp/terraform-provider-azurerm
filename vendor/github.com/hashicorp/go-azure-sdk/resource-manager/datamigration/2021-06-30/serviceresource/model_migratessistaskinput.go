package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSsisTaskInput struct {
	SourceConnectionInfo SqlConnectionInfo `json:"sourceConnectionInfo"`
	SsisMigrationInfo    SsisMigrationInfo `json:"ssisMigrationInfo"`
	TargetConnectionInfo SqlConnectionInfo `json:"targetConnectionInfo"`
}
