package dataflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerQuerySource struct {
	Dataset             *DatasetReference       `json:"dataset,omitempty"`
	Description         *string                 `json:"description,omitempty"`
	Flowlet             *DataFlowReference      `json:"flowlet,omitempty"`
	LinkedService       *LinkedServiceReference `json:"linkedService,omitempty"`
	Name                string                  `json:"name"`
	SchemaLinkedService *LinkedServiceReference `json:"schemaLinkedService,omitempty"`
	Script              *string                 `json:"script,omitempty"`
}
