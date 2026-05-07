package rbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlRoleAssignmentResource struct {
	PrincipalId      *string `json:"principalId,omitempty"`
	RoleDefinitionId *string `json:"roleDefinitionId,omitempty"`
	Scope            *string `json:"scope,omitempty"`
}
