package dataflow

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowBuiltInTransformationSettings struct {
	Datasets            *[]DataflowBuiltInTransformationDataset `json:"datasets,omitempty"`
	Filter              *[]DataflowBuiltInTransformationFilter  `json:"filter,omitempty"`
	Map                 *[]DataflowBuiltInTransformationMap     `json:"map,omitempty"`
	SchemaRef           *string                                 `json:"schemaRef,omitempty"`
	SerializationFormat *TransformationSerializationFormat      `json:"serializationFormat,omitempty"`
}
