package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUserTablesSqlSyncTaskOutput struct {
	DatabasesToSourceTables *map[string][]DatabaseTable `json:"databasesToSourceTables,omitempty"`
	DatabasesToTargetTables *map[string][]DatabaseTable `json:"databasesToTargetTables,omitempty"`
	TableValidationErrors   *map[string][]string        `json:"tableValidationErrors,omitempty"`
	ValidationErrors        *[]ReportableException      `json:"validationErrors,omitempty"`
}
