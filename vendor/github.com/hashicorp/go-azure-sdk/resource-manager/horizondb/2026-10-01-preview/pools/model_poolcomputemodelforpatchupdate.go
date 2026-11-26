package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolComputeModelForPatchUpdate struct {
	ComputeGeneration *int64                `json:"computeGeneration,omitempty"`
	ComputeProcessor  *ComputeProcessorType `json:"computeProcessor,omitempty"`
	Type              *PoolComputeModelType `json:"type,omitempty"`
	VCores            *int64                `json:"vCores,omitempty"`
}
