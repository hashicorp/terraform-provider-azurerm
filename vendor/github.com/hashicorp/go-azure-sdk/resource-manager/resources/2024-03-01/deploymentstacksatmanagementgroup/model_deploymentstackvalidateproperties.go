package deploymentstacksatmanagementgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStackValidateProperties struct {
	ActionOnUnmanage   *ActionOnUnmanage               `json:"actionOnUnmanage,omitempty"`
	CorrelationId      *string                         `json:"correlationId,omitempty"`
	DenySettings       *DenySettings                   `json:"denySettings,omitempty"`
	DeploymentScope    *string                         `json:"deploymentScope,omitempty"`
	Description        *string                         `json:"description,omitempty"`
	Parameters         *map[string]DeploymentParameter `json:"parameters,omitempty"`
	TemplateLink       *DeploymentStacksTemplateLink   `json:"templateLink,omitempty"`
	ValidatedResources *[]ResourceReference            `json:"validatedResources,omitempty"`
}
