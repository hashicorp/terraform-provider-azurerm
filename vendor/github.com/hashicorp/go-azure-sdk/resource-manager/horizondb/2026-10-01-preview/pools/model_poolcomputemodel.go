package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolComputeModel struct {
	ComputeGeneration int64                `json:"computeGeneration"`
	ComputeProcessor  ComputeProcessorType `json:"computeProcessor"`
	Type              PoolComputeModelType `json:"type"`
	VCores            int64                `json:"vCores"`
}
