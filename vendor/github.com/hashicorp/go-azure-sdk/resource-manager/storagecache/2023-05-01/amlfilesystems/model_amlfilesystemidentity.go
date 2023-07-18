package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemIdentity struct {
	PrincipalId            *string                                      `json:"principalId,omitempty"`
	TenantId               *string                                      `json:"tenantId,omitempty"`
	Type                   *AmlFilesystemIdentityType                   `json:"type,omitempty"`
	UserAssignedIdentities *map[string]UserAssignedIdentitiesProperties `json:"userAssignedIdentities,omitempty"`
}
