package collectorpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectorPolicyPropertiesFormat struct {
	EmissionPolicies  *[]EmissionPoliciesPropertiesFormat `json:"emissionPolicies,omitempty"`
	IngestionPolicy   *IngestionPolicyPropertiesFormat    `json:"ingestionPolicy,omitempty"`
	ProvisioningState *ProvisioningState                  `json:"provisioningState,omitempty"`
}
