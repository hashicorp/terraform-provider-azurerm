package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances"
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

func RoleAssignmentScheduleID(input string) (*string, error) {
	re := regexp.MustCompile(`^.+/providers/Microsoft.Authorization/roleEligibilitySchedules/(.+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 2 {
		return nil, fmt.Errorf("parsing %s", input)
	}
	return &matches[1], nil
}

func RoleAssignmentScheduleIdFromInstance(r *roleassignmentscheduleinstances.RoleAssignmentScheduleInstance) (*string, error) {
	re := regexp.MustCompile(`^.+/providers/Microsoft.Authorization/roleAssignmentSchedules/(.+)`)
	matches := re.FindStringSubmatch(*r.Properties.RoleAssignmentScheduleId)
	if len(matches) != 2 {
		return nil, fmt.Errorf("parsing %s", *r.Properties.RoleAssignmentScheduleId)
	}
	return &matches[1], nil
}

func RoleEligibilityScheduleIdFromInstance(r *roleeligibilityscheduleinstances.RoleEligibilityScheduleInstance) (*string, error) {
	re := regexp.MustCompile(`^.+/providers/Microsoft.Authorization/roleEligibilitySchedules/(.+)`)
	matches := re.FindStringSubmatch(*r.Properties.RoleEligibilityScheduleId)
	if len(matches) != 2 {
		return nil, fmt.Errorf("parsing %s", *r.Properties.RoleEligibilityScheduleId)
	}
	return &matches[1], nil
}
