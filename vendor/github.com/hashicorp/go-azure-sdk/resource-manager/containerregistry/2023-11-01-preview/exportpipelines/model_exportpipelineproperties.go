package exportpipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportPipelineProperties struct {
	Options           *[]PipelineOptions             `json:"options,omitempty"`
	ProvisioningState *ProvisioningState             `json:"provisioningState,omitempty"`
	Target            ExportPipelineTargetProperties `json:"target"`
}
