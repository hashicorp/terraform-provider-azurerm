package deploymentscripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedServiceIdentity struct {
	TenantId               *string                          `json:"tenantId,omitempty"`
	Type                   *ManagedServiceIdentityType      `json:"type,omitempty"`
	UserAssignedIdentities *map[string]UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}
