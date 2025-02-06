package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CsmDeploymentStatusProperties struct {
	DeploymentId                *string                `json:"deploymentId,omitempty"`
	Errors                      *[]ErrorEntity         `json:"errors,omitempty"`
	FailedInstancesLogs         *[]string              `json:"failedInstancesLogs,omitempty"`
	NumberOfInstancesFailed     *int64                 `json:"numberOfInstancesFailed,omitempty"`
	NumberOfInstancesInProgress *int64                 `json:"numberOfInstancesInProgress,omitempty"`
	NumberOfInstancesSuccessful *int64                 `json:"numberOfInstancesSuccessful,omitempty"`
	Status                      *DeploymentBuildStatus `json:"status,omitempty"`
}
