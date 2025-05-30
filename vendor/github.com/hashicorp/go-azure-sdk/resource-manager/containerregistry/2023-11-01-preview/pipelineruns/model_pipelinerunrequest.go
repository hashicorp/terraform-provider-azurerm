package pipelineruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineRunRequest struct {
	Artifacts          *[]string                    `json:"artifacts,omitempty"`
	CatalogDigest      *string                      `json:"catalogDigest,omitempty"`
	PipelineResourceId *string                      `json:"pipelineResourceId,omitempty"`
	Source             *PipelineRunSourceProperties `json:"source,omitempty"`
	Target             *PipelineRunTargetProperties `json:"target,omitempty"`
}
