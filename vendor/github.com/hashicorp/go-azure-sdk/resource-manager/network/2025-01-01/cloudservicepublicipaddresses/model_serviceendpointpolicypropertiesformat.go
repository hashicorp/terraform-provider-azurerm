package cloudservicepublicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpointPolicyPropertiesFormat struct {
	ContextualServiceEndpointPolicies *[]string                          `json:"contextualServiceEndpointPolicies,omitempty"`
	ProvisioningState                 *ProvisioningState                 `json:"provisioningState,omitempty"`
	ResourceGuid                      *string                            `json:"resourceGuid,omitempty"`
	ServiceAlias                      *string                            `json:"serviceAlias,omitempty"`
	ServiceEndpointPolicyDefinitions  *[]ServiceEndpointPolicyDefinition `json:"serviceEndpointPolicyDefinitions,omitempty"`
	Subnets                           *[]Subnet                          `json:"subnets,omitempty"`
}
