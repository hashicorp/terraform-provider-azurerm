package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedIdentity struct {
	TenantId               *string                  `json:"tenantId,omitempty"`
	Type                   IdentityType             `json:"type"`
	UserAssignedIdentities *map[string]UserIdentity `json:"userAssignedIdentities,omitempty"`
}
