package clusterprincipalassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPrincipalProperties struct {
	AadObjectId       *string              `json:"aadObjectId,omitempty"`
	PrincipalId       string               `json:"principalId"`
	PrincipalName     *string              `json:"principalName,omitempty"`
	PrincipalType     PrincipalType        `json:"principalType"`
	ProvisioningState *ProvisioningState   `json:"provisioningState,omitempty"`
	Role              ClusterPrincipalRole `json:"role"`
	TenantId          *string              `json:"tenantId,omitempty"`
	TenantName        *string              `json:"tenantName,omitempty"`
}
