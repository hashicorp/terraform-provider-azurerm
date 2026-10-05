package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProperties struct {
	PrincipalId            *string                          `json:"principalId,omitempty"`
	TenantId               *string                          `json:"tenantId,omitempty"`
	Type                   ManagedIdentityType              `json:"type"`
	UserAssignedIdentities *map[string]UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}
