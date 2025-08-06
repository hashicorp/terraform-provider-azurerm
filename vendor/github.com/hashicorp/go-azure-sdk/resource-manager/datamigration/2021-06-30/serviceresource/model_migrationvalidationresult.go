package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationValidationResult struct {
	Id             *string                                              `json:"id,omitempty"`
	MigrationId    *string                                              `json:"migrationId,omitempty"`
	Status         *ValidationStatus                                    `json:"status,omitempty"`
	SummaryResults *map[string]MigrationValidationDatabaseSummaryResult `json:"summaryResults,omitempty"`
}
