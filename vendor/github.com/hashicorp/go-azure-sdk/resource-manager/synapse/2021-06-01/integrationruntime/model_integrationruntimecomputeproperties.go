package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeComputeProperties struct {
	DataFlowProperties           *IntegrationRuntimeDataFlowProperties `json:"dataFlowProperties,omitempty"`
	Location                     *string                               `json:"location,omitempty"`
	MaxParallelExecutionsPerNode *int64                                `json:"maxParallelExecutionsPerNode,omitempty"`
	NodeSize                     *string                               `json:"nodeSize,omitempty"`
	NumberOfNodes                *int64                                `json:"numberOfNodes,omitempty"`
	VNetProperties               *IntegrationRuntimeVNetProperties     `json:"vNetProperties,omitempty"`
}
