// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RoleManagementPolicyId struct {
	RoleDefinitionId string
	Scope            string
}

var _ resourceids.Id = RoleManagementPolicyId{}

func NewRoleManagementPolicyId(roleDefinitionId string, scope string) RoleManagementPolicyId {
	return RoleManagementPolicyId{
		RoleDefinitionId: roleDefinitionId,
		Scope:            scope,
	}
}

// RoleManagementPolicyID parses 'input' into a RoleManagementPolicyId
func RoleManagementPolicyID(input string) (*RoleManagementPolicyId, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Role Management Policy ID, invalid format %q", input)
	}

	return &RoleManagementPolicyId{
		RoleDefinitionId: parts[0],
		Scope:            parts[1],
	}, nil
}

func (id RoleManagementPolicyId) ID() string {
	return fmt.Sprintf("%s|%s", id.RoleDefinitionId, id.Scope)
}

func (id RoleManagementPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Role Definition ID: %q", id.RoleDefinitionId),
	}
	if id.Scope != "" {
		components = append(components, fmt.Sprintf("Scope: %q", id.Scope))
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}

func ValidateRoleManagementPolicyId(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := RoleManagementPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
