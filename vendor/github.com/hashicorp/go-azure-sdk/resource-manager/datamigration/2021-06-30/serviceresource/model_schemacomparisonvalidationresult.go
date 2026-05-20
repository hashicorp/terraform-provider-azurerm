package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaComparisonValidationResult struct {
	SchemaDifferences         *SchemaComparisonValidationResultType `json:"schemaDifferences,omitempty"`
	SourceDatabaseObjectCount *map[string]int64                     `json:"sourceDatabaseObjectCount,omitempty"`
	TargetDatabaseObjectCount *map[string]int64                     `json:"targetDatabaseObjectCount,omitempty"`
	ValidationErrors          *ValidationError                      `json:"validationErrors,omitempty"`
}
