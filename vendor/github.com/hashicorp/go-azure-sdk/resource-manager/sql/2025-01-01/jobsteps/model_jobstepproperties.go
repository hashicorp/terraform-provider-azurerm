package jobsteps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStepProperties struct {
	Action           JobStepAction            `json:"action"`
	Credential       *string                  `json:"credential,omitempty"`
	ExecutionOptions *JobStepExecutionOptions `json:"executionOptions,omitempty"`
	Output           *JobStepOutput           `json:"output,omitempty"`
	StepId           *int64                   `json:"stepId,omitempty"`
	TargetGroup      string                   `json:"targetGroup"`
}
