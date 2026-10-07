package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineReference struct {
	Name          *string               `json:"name,omitempty"`
	ReferenceName string                `json:"referenceName"`
	Type          PipelineReferenceType `json:"type"`
}
