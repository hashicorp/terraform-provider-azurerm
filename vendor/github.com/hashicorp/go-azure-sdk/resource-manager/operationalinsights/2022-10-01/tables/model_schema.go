package tables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schema struct {
	Categories      *[]string         `json:"categories,omitempty"`
	Columns         *[]Column         `json:"columns,omitempty"`
	Description     *string           `json:"description,omitempty"`
	DisplayName     *string           `json:"displayName,omitempty"`
	Labels          *[]string         `json:"labels,omitempty"`
	Name            *string           `json:"name,omitempty"`
	Solutions       *[]string         `json:"solutions,omitempty"`
	Source          *SourceEnum       `json:"source,omitempty"`
	StandardColumns *[]Column         `json:"standardColumns,omitempty"`
	TableSubType    *TableSubTypeEnum `json:"tableSubType,omitempty"`
	TableType       *TableTypeEnum    `json:"tableType,omitempty"`
}
