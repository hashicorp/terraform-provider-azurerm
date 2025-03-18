package pipelineruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineRunSourceProperties struct {
	Name *string                `json:"name,omitempty"`
	Type *PipelineRunSourceType `json:"type,omitempty"`
}
