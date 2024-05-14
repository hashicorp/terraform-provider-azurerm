// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedules"
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
	return fmt.Sprintf("%s: (%s)", "Role Management Policy", segmentsStr)
}

func PimRoleAssignmentID(input string) (*PimRoleAssignmentId, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("could not parse Role Management Policy ID, invalid format %q", input)
	}

	pimRoleAssignmentId := PimRoleAssignmentId{
		Scope:            parts[0],
		RoleDefinitionId: parts[1],
		PrincipalId:      parts[2],
	}

	return &pimRoleAssignmentId, nil
}

func (id PimRoleAssignmentId) ScopeID() commonids.ScopeId {
	return commonids.NewScopeID(id.Scope)
}

func RoleEligibilityScheduleRequestIdFromSchedule(r *roleeligibilityschedules.RoleEligibilitySchedule) (*string, error) {
	re := regexp.MustCompile(`^.+/providers/Microsoft.Authorization/roleEligibilityScheduleRequests/(.+)`)
	matches := re.FindStringSubmatch(*r.Properties.RoleEligibilityScheduleRequestId)
	if len(matches) != 2 {
		return nil, fmt.Errorf("parsing %s", *r.Properties.RoleEligibilityScheduleRequestId)
	}
	return &matches[1], nil
}
