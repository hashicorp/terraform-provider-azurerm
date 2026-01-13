package integrationruntimedisableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeComputeProperties struct {
	CopyComputeScaleProperties             *CopyComputeScaleProperties             `json:"copyComputeScaleProperties,omitempty"`
	DataFlowProperties                     *IntegrationRuntimeDataFlowProperties   `json:"dataFlowProperties,omitempty"`
	Location                               *string                                 `json:"location,omitempty"`
	MaxParallelExecutionsPerNode           *int64                                  `json:"maxParallelExecutionsPerNode,omitempty"`
	NodeSize                               *string                                 `json:"nodeSize,omitempty"`
	NumberOfNodes                          *int64                                  `json:"numberOfNodes,omitempty"`
	PipelineExternalComputeScaleProperties *PipelineExternalComputeScaleProperties `json:"pipelineExternalComputeScaleProperties,omitempty"`
	VNetProperties                         *IntegrationRuntimeVNetProperties       `json:"vNetProperties,omitempty"`
}
