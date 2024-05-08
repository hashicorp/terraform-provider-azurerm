package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RoleManagementPolicyId struct {
	Scope string
	Name  string
}

var _ resourceids.Id = RoleManagementPolicyId{}

func NewRoleManagementPolicyID(scope, name string) *RoleManagementPolicyId {
	return &RoleManagementPolicyId{
		Scope: scope,
		Name:  name,
	}
}

func RoleManagementPolicyID(id string) (*RoleManagementPolicyId, error) {
	parts := strings.Split(id, "/providers/Microsoft.Authorization/roleManagementPolicies/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid role management policy ID: %s", id)
	}

	return NewRoleManagementPolicyID(parts[0], parts[1]), nil
}

func (id RoleManagementPolicyId) ID() string {
	fmtString := "%s/providers/Microsoft.Authorization/roleManagementPolicies/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.Name)
}

func (id RoleManagementPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Policy Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Role Management Policy ID", segmentsStr)
}
