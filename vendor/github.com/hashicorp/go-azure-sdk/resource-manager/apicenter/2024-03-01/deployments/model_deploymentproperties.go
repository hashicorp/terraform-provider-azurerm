package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentProperties struct {
	CustomProperties *interface{}      `json:"customProperties,omitempty"`
	DefinitionId     *string           `json:"definitionId,omitempty"`
	Description      *string           `json:"description,omitempty"`
	EnvironmentId    *string           `json:"environmentId,omitempty"`
	Server           *DeploymentServer `json:"server,omitempty"`
	State            *DeploymentState  `json:"state,omitempty"`
	Title            *string           `json:"title,omitempty"`
}
