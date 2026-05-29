package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticDataTableResponseColumn struct {
	ColumnName *string `json:"columnName,omitempty"`
	ColumnType *string `json:"columnType,omitempty"`
	DataType   *string `json:"dataType,omitempty"`
}
