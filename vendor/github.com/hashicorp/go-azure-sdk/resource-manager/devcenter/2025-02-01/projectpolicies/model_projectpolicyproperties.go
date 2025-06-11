package projectpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectPolicyProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ResourcePolicies  *[]ResourcePolicy  `json:"resourcePolicies,omitempty"`
	Scopes            *[]string          `json:"scopes,omitempty"`
}
