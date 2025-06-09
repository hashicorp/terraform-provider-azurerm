package dataflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataFlowReference struct {
	DatasetParameters *interface{}            `json:"datasetParameters,omitempty"`
	Parameters        *map[string]interface{} `json:"parameters,omitempty"`
	ReferenceName     string                  `json:"referenceName"`
	Type              DataFlowReferenceType   `json:"type"`
}
