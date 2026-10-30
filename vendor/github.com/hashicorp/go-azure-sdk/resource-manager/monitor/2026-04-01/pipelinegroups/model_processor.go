package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Processor struct {
	Batch             *BatchProcessor             `json:"batch,omitempty"`
	Name              string                      `json:"name"`
	TransformLanguage *TransformLanguageProcessor `json:"transformLanguage,omitempty"`
	Type              ProcessorType               `json:"type"`
}
