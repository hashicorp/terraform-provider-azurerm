package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolPropertiesUpdate struct {
	AutoScaling             *NodePoolAutoScaling          `json:"autoScaling,omitempty"`
	Labels                  *[]Label                      `json:"labels,omitempty"`
	NodeDrainTimeoutMinutes *int64                        `json:"nodeDrainTimeoutMinutes,omitempty"`
	Replicas                *int64                        `json:"replicas,omitempty"`
	Taints                  *[]Taint                      `json:"taints,omitempty"`
	Version                 *NodePoolVersionProfileUpdate `json:"version,omitempty"`
}
