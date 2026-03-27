package dataflow

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowBuiltInTransformationDataset struct {
	Description *string  `json:"description,omitempty"`
	Expression  *string  `json:"expression,omitempty"`
	Inputs      []string `json:"inputs"`
	Key         string   `json:"key"`
	SchemaRef   *string  `json:"schemaRef,omitempty"`
}
