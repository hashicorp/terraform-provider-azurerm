package goldengateconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentConnectionAssignmentProperties struct {
	AliasName         *string                                       `json:"aliasName,omitempty"`
	CompartmentId     *string                                       `json:"compartmentId,omitempty"`
	ConnectionId      string                                        `json:"connectionId"`
	ConnectionName    *string                                       `json:"connectionName,omitempty"`
	DeploymentId      string                                        `json:"deploymentId"`
	DeploymentName    *string                                       `json:"deploymentName,omitempty"`
	LifecycleState    *GoldenGateConnectionAssignmentLifecycleState `json:"lifecycleState,omitempty"`
	Ocid              *string                                       `json:"ocid,omitempty"`
	ProvisioningState *AzureResourceProvisioningState               `json:"provisioningState,omitempty"`
	TimeCreated       *string                                       `json:"timeCreated,omitempty"`
	TimeUpdated       *string                                       `json:"timeUpdated,omitempty"`
}
