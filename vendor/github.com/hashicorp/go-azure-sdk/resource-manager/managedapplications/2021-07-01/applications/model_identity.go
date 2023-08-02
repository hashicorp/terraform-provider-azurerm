package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Identity struct {
	PrincipalId            *string                                  `json:"principalId,omitempty"`
	TenantId               *string                                  `json:"tenantId,omitempty"`
	Type                   *ResourceIdentityType                    `json:"type,omitempty"`
	UserAssignedIdentities *map[string]UserAssignedResourceIdentity `json:"userAssignedIdentities,omitempty"`
}
