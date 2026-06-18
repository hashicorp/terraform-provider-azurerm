package roleassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentDetails struct {
	Id               *string `json:"id,omitempty"`
	PrincipalId      *string `json:"principalId,omitempty"`
	PrincipalType    *string `json:"principalType,omitempty"`
	RoleDefinitionId *string `json:"roleDefinitionId,omitempty"`
	Scope            *string `json:"scope,omitempty"`
}
