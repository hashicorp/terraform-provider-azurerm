package integrationserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationServiceEnvironmentProperties struct {
	EncryptionConfiguration         *IntegrationServiceEnvironmenEncryptionConfiguration `json:"encryptionConfiguration,omitempty"`
	EndpointsConfiguration          *FlowEndpointsConfiguration                          `json:"endpointsConfiguration,omitempty"`
	IntegrationServiceEnvironmentId *string                                              `json:"integrationServiceEnvironmentId,omitempty"`
	NetworkConfiguration            *NetworkConfiguration                                `json:"networkConfiguration,omitempty"`
	ProvisioningState               *WorkflowProvisioningState                           `json:"provisioningState,omitempty"`
	State                           *WorkflowState                                       `json:"state,omitempty"`
}
