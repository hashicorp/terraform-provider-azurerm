package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePurgeBodyFilters struct {
	Column   *string      `json:"column,omitempty"`
	Key      *string      `json:"key,omitempty"`
	Operator *string      `json:"operator,omitempty"`
	Value    *interface{} `json:"value,omitempty"`
}
