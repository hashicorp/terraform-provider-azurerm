package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualApplianceMigrationStatus struct {
	MigrationPhase       *string        `json:"migrationPhase,omitempty"`
	MigrationPhaseStatus *string        `json:"migrationPhaseStatus,omitempty"`
	MigrationType        *MigrationType `json:"migrationType,omitempty"`
}
