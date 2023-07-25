package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArtifactDeploymentStatusProperties struct {
	ArtifactsApplied *int64  `json:"artifactsApplied,omitempty"`
	DeploymentStatus *string `json:"deploymentStatus,omitempty"`
	TotalArtifacts   *int64  `json:"totalArtifacts,omitempty"`
}
