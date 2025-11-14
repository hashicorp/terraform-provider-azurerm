package dataflow

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowOperation struct {
	BuiltInTransformationSettings *DataflowBuiltInTransformationSettings `json:"builtInTransformationSettings,omitempty"`
	DestinationSettings           *DataflowDestinationOperationSettings  `json:"destinationSettings,omitempty"`
	Name                          *string                                `json:"name,omitempty"`
	OperationType                 OperationType                          `json:"operationType"`
	SourceSettings                *DataflowSourceOperationSettings       `json:"sourceSettings,omitempty"`
}
