package identities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedIdentityProperties struct {
	ClientId       *string         `json:"clientId,omitempty"`
	IsolationScope *IsolationScope `json:"isolationScope,omitempty"`
	PrincipalId    *string         `json:"principalId,omitempty"`
	TenantId       *string         `json:"tenantId,omitempty"`
}
