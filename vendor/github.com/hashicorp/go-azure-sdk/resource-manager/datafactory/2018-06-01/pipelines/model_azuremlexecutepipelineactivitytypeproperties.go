package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLExecutePipelineActivityTypeProperties struct {
	ContinueOnStepFailure *bool              `json:"continueOnStepFailure,omitempty"`
	DataPathAssignments   *interface{}       `json:"dataPathAssignments,omitempty"`
	ExperimentName        *interface{}       `json:"experimentName,omitempty"`
	MlParentRunId         *interface{}       `json:"mlParentRunId,omitempty"`
	MlPipelineEndpointId  *interface{}       `json:"mlPipelineEndpointId,omitempty"`
	MlPipelineId          *interface{}       `json:"mlPipelineId,omitempty"`
	MlPipelineParameters  *map[string]string `json:"mlPipelineParameters,omitempty"`
	Version               *interface{}       `json:"version,omitempty"`
}
