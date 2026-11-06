package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolProperties struct {
	AutoRepair              *bool                   `json:"autoRepair,omitempty"`
	AutoScaling             *NodePoolAutoScaling    `json:"autoScaling,omitempty"`
	Labels                  *[]Label                `json:"labels,omitempty"`
	NodeDrainTimeoutMinutes *int64                  `json:"nodeDrainTimeoutMinutes,omitempty"`
	Platform                NodePoolPlatformProfile `json:"platform"`
	ProvisioningState       *ProvisioningState      `json:"provisioningState,omitempty"`
	Replicas                *int64                  `json:"replicas,omitempty"`
	Status                  *NodePoolResourceStatus `json:"status,omitempty"`
	Taints                  *[]Taint                `json:"taints,omitempty"`
	Version                 *NodePoolVersionProfile `json:"version,omitempty"`
}
