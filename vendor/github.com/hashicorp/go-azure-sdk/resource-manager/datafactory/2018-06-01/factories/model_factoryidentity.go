package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryIdentity struct {
	PrincipalId            *string                 `json:"principalId,omitempty"`
	TenantId               *string                 `json:"tenantId,omitempty"`
	Type                   FactoryIdentityType     `json:"type"`
	UserAssignedIdentities *map[string]interface{} `json:"userAssignedIdentities,omitempty"`
}
