package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
)

type RoleManagementPolicyID struct {
	Scope                    string
	RoleManagementPolicyName string
	RoleDefinitionId         string
}

func NewRoleManagementPolicyID(scope string, roleManagementPolicyName string, roleDefinitionId string) RoleManagementPolicyID {
	return RoleManagementPolicyID{
		Scope:                    scope,
		RoleManagementPolicyName: roleManagementPolicyName,
		RoleDefinitionId:         roleDefinitionId,
	}
}

func (id RoleManagementPolicyID) ID() string {
	fmtString := "%s/providers/Microsoft.Authorization/roleManagementPolicies/%s|%s"
	return fmt.Sprintf(fmtString, id.Scope, id.RoleManagementPolicyName, id.RoleDefinitionId)
}

func (id RoleManagementPolicyID) ScopedRoleManagementPolicyId() rolemanagementpolicies.ScopedRoleManagementPolicyId {
	return rolemanagementpolicies.NewScopedRoleManagementPolicyID(id.Scope, id.RoleManagementPolicyName)
}

func (id RoleManagementPolicyID) String() string {
	segments := []string{
		fmt.Sprintf("RoleManagementPolicyName %q", id.RoleManagementPolicyName),
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Role Definition Id %q", id.RoleDefinitionId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Role Management Policy", segmentsStr)
}

// RoleManagementPolicyId is a pseudo ID for storing Role Definition ID parameter as this it not retrievable from API
// It is formed of the Azure Resource ID for the Role Management Policy ID and the Role Definition ID it is created against
func RoleManagementPolicyId(input string) (*RoleManagementPolicyID, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Role Management Policy ID, invalid format %q", input)
	}

	roleManagementPolicyID := RoleManagementPolicyID{}

	rawRoleManagementPolicyId := parts[0]
	rawRoleDefinitionId := parts[1]

	roleManagementPolicyId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(rawRoleManagementPolicyId)
	if err != nil {
		return nil, err
	}
	roleManagementPolicyID.Scope = roleManagementPolicyId.Scope
	roleManagementPolicyID.RoleManagementPolicyName = roleManagementPolicyId.RoleManagementPolicyName

	roleManagementPolicyID.RoleDefinitionId = rawRoleDefinitionId

	return &roleManagementPolicyID, nil
}
