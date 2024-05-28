// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"
)

type PimRoleAssignmentId struct {
	Scope            string
	RoleDefinitionId string
	PrincipalId      string
}

func NewPimRoleAssignmentID(scope string, roleDefinitionId string, principalId string) PimRoleAssignmentId {
	return PimRoleAssignmentId{
		Scope:            scope,
		RoleDefinitionId: roleDefinitionId,
		PrincipalId:      principalId,
	}
}

func (id PimRoleAssignmentId) ID() string {
	fmtString := "%s|%s|%s"
	return fmt.Sprintf(fmtString, id.Scope, id.RoleDefinitionId, id.PrincipalId)
}

func (id PimRoleAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Principal Id %q", id.PrincipalId),
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Role Definition Id %q", id.RoleDefinitionId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "PIM Role Assignment", segmentsStr)
}

func PimRoleAssignmentID(input string) (*PimRoleAssignmentId, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("could not parse PIM Role Assignment ID, invalid format %q", input)
	}

	pimRoleAssignmentId := PimRoleAssignmentId{
		Scope:            parts[0],
		RoleDefinitionId: parts[1],
		PrincipalId:      parts[2],
	}

	return &pimRoleAssignmentId, nil
}
