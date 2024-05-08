package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterResourceProperties struct {
	ClusterProfile    ClusterProfile      `json:"clusterProfile"`
	ClusterType       string              `json:"clusterType"`
	ComputeProfile    ComputeProfile      `json:"computeProfile"`
	DeploymentId      *string             `json:"deploymentId,omitempty"`
	ProvisioningState *ProvisioningStatus `json:"provisioningState,omitempty"`
	Status            *string             `json:"status,omitempty"`
}
