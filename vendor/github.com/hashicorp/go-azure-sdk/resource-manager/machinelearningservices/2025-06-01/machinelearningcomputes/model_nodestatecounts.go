package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeStateCounts struct {
	IdleNodeCount      *int64 `json:"idleNodeCount,omitempty"`
	LeavingNodeCount   *int64 `json:"leavingNodeCount,omitempty"`
	PreemptedNodeCount *int64 `json:"preemptedNodeCount,omitempty"`
	PreparingNodeCount *int64 `json:"preparingNodeCount,omitempty"`
	RunningNodeCount   *int64 `json:"runningNodeCount,omitempty"`
	UnusableNodeCount  *int64 `json:"unusableNodeCount,omitempty"`
}
