package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DWCopyCommandDefaultValue struct {
	ColumnName   *interface{} `json:"columnName,omitempty"`
	DefaultValue *interface{} `json:"defaultValue,omitempty"`
}
