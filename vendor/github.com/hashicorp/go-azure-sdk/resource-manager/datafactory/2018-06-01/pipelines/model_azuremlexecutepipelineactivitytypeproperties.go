package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLExecutePipelineActivityTypeProperties struct {
	ContinueOnStepFailure *bool              `json:"continueOnStepFailure,omitempty"`
	DataPathAssignments   *interface{}       `json:"dataPathAssignments,omitempty"`
	ExperimentName        *string            `json:"experimentName,omitempty"`
	MlParentRunId         *string            `json:"mlParentRunId,omitempty"`
	MlPipelineEndpointId  *string            `json:"mlPipelineEndpointId,omitempty"`
	MlPipelineId          *string            `json:"mlPipelineId,omitempty"`
	MlPipelineParameters  *map[string]string `json:"mlPipelineParameters,omitempty"`
	Version               *string            `json:"version,omitempty"`
}
