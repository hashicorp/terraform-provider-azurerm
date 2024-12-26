package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FixedScaleSettings struct {
	NodeDeallocationOption *ComputeNodeDeallocationOption `json:"nodeDeallocationOption,omitempty"`
	ResizeTimeout          *string                        `json:"resizeTimeout,omitempty"`
	TargetDedicatedNodes   *int64                         `json:"targetDedicatedNodes,omitempty"`
	TargetLowPriorityNodes *int64                         `json:"targetLowPriorityNodes,omitempty"`
}
