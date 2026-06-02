package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnErrorDeploymentExtended struct {
	DeploymentName    *string                `json:"deploymentName,omitempty"`
	ProvisioningState *string                `json:"provisioningState,omitempty"`
	Type              *OnErrorDeploymentType `json:"type,omitempty"`
}
