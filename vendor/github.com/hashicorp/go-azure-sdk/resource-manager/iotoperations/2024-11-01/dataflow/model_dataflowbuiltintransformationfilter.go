package dataflow

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowBuiltInTransformationFilter struct {
	Description *string     `json:"description,omitempty"`
	Expression  string      `json:"expression"`
	Inputs      []string    `json:"inputs"`
	Type        *FilterType `json:"type,omitempty"`
}
