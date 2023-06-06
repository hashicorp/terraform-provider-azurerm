package jobdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobDefinitionProperties struct {
	AgentName              *string            `json:"agentName,omitempty"`
	AgentResourceId        *string            `json:"agentResourceId,omitempty"`
	CopyMode               CopyMode           `json:"copyMode"`
	Description            *string            `json:"description,omitempty"`
	LatestJobRunName       *string            `json:"latestJobRunName,omitempty"`
	LatestJobRunResourceId *string            `json:"latestJobRunResourceId,omitempty"`
	LatestJobRunStatus     *JobRunStatus      `json:"latestJobRunStatus,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	SourceName             string             `json:"sourceName"`
	SourceResourceId       *string            `json:"sourceResourceId,omitempty"`
	SourceSubpath          *string            `json:"sourceSubpath,omitempty"`
	TargetName             string             `json:"targetName"`
	TargetResourceId       *string            `json:"targetResourceId,omitempty"`
	TargetSubpath          *string            `json:"targetSubpath,omitempty"`
}
