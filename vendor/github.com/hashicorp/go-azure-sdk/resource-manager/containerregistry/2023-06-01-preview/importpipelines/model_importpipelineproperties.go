package importpipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportPipelineProperties struct {
	Options           *[]PipelineOptions             `json:"options,omitempty"`
	ProvisioningState *ProvisioningState             `json:"provisioningState,omitempty"`
	Source            ImportPipelineSourceProperties `json:"source"`
	Trigger           *PipelineTriggerProperties     `json:"trigger,omitempty"`
}
