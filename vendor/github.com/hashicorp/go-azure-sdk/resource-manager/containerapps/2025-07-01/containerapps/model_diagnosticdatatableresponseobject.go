package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticDataTableResponseObject struct {
	Columns   *[]DiagnosticDataTableResponseColumn `json:"columns,omitempty"`
	Rows      *[]interface{}                       `json:"rows,omitempty"`
	TableName *string                              `json:"tableName,omitempty"`
}
