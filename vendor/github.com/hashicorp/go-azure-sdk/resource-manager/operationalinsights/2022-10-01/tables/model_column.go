package tables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Column struct {
	DataTypeHint     *ColumnDataTypeHintEnum `json:"dataTypeHint,omitempty"`
	Description      *string                 `json:"description,omitempty"`
	DisplayName      *string                 `json:"displayName,omitempty"`
	IsDefaultDisplay *bool                   `json:"isDefaultDisplay,omitempty"`
	IsHidden         *bool                   `json:"isHidden,omitempty"`
	Name             *string                 `json:"name,omitempty"`
	Type             *ColumnTypeEnum         `json:"type,omitempty"`
}
