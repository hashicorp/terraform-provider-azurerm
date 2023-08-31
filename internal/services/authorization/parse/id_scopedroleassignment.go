package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
)

type ScopedRoleAssignmentId struct {
	ScopedId roleassignments.ScopedRoleAssignmentId
	TenantId string
}

func NewScopedRoleAssignmentID(scope string, roleAssignmentName string, tenantId string) ScopedRoleAssignmentId {
	return ScopedRoleAssignmentId{
		ScopedId: roleassignments.NewScopedRoleAssignmentID(scope, roleAssignmentName),
		TenantId: tenantId,
	}
}

func ScopedRoleAssignmentID(input string) (*ScopedRoleAssignmentId, error) {
	azureResourceId, tenantId := DestructRoleAssignmentId(input)
	scopedId, err := roleassignments.ParseScopedRoleAssignmentID(azureResourceId)
	if err != nil {
		return nil, err
	}

	return &ScopedRoleAssignmentId{ScopedId: *scopedId, TenantId: tenantId}, nil
}

func ValidateScopedRoleAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ScopedRoleAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func (id ScopedRoleAssignmentId) ID() string {
	return ConstructRoleAssignmentId(id.ScopedId.ID(), id.TenantId)
}

func (id ScopedRoleAssignmentId) String() string {
	components := []string{
		id.ScopedId.String(),
	}

	if id.TenantId != "" {
		components = append(components, fmt.Sprintf("Tenant ID: %s", id.TenantId))
	}

	return fmt.Sprintf("Scoped Role Assignment (%s)", strings.Join(components, "\n"))
}
