package integrationruntimeenableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineExternalComputeScaleProperties struct {
	NumberOfExternalNodes *int64 `json:"numberOfExternalNodes,omitempty"`
	NumberOfPipelineNodes *int64 `json:"numberOfPipelineNodes,omitempty"`
	TimeToLive            *int64 `json:"timeToLive,omitempty"`
}
