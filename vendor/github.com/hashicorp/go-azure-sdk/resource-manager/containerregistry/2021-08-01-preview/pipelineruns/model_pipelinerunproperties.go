package pipelineruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineRunProperties struct {
	ForceUpdateTag    *string              `json:"forceUpdateTag,omitempty"`
	ProvisioningState *ProvisioningState   `json:"provisioningState,omitempty"`
	Request           *PipelineRunRequest  `json:"request,omitempty"`
	Response          *PipelineRunResponse `json:"response,omitempty"`
}
