package dataflow

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowSourceOperationSettings struct {
	AssetRef            *string                    `json:"assetRef,omitempty"`
	DataSources         []string                   `json:"dataSources"`
	EndpointRef         string                     `json:"endpointRef"`
	SchemaRef           *string                    `json:"schemaRef,omitempty"`
	SerializationFormat *SourceSerializationFormat `json:"serializationFormat,omitempty"`
}
