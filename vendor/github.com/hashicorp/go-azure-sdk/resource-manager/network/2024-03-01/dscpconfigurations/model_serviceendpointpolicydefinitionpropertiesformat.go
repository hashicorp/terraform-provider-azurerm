package dscpconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpointPolicyDefinitionPropertiesFormat struct {
	Description       *string            `json:"description,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Service           *string            `json:"service,omitempty"`
	ServiceResources  *[]string          `json:"serviceResources,omitempty"`
}
