package accountmigrations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountMigrationProperties struct {
	MigrationFailedDetailedReason *string          `json:"migrationFailedDetailedReason,omitempty"`
	MigrationFailedReason         *string          `json:"migrationFailedReason,omitempty"`
	MigrationStatus               *MigrationStatus `json:"migrationStatus,omitempty"`
	TargetSkuName                 SkuName          `json:"targetSkuName"`
}
