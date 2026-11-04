package jobagents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobAgentIdentity struct {
	TenantId               *string                                  `json:"tenantId,omitempty"`
	Type                   JobAgentIdentityType                     `json:"type"`
	UserAssignedIdentities *map[string]JobAgentUserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}
