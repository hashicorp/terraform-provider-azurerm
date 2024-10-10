package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUserTablesOracleTaskOutput struct {
	SchemaName       *string                `json:"schemaName,omitempty"`
	Tables           *[]DatabaseTable       `json:"tables,omitempty"`
	ValidationErrors *[]ReportableException `json:"validationErrors,omitempty"`
}
