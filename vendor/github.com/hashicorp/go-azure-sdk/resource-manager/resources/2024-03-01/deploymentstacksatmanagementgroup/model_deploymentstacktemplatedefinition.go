package deploymentstacksatmanagementgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStackTemplateDefinition struct {
	Template     *map[string]interface{}       `json:"template,omitempty"`
	TemplateLink *DeploymentStacksTemplateLink `json:"templateLink,omitempty"`
}
