package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlComputeNodesInformation struct {
	NextLink *string                      `json:"nextLink,omitempty"`
	Nodes    *[]AmlComputeNodeInformation `json:"nodes,omitempty"`
}
