// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Resource Manager Role Definition ID
type RoleDefinitionResourceId struct {
	RoleDefinitionId string
	Scope            string
}

var _ resourceids.Id = RoleDefinitionResourceId{}

func NewRoleDefinitionResourceId(roleDefinitionId, scope string) RoleDefinitionResourceId {
	return RoleDefinitionResourceId{
		RoleDefinitionId: roleDefinitionId,
		Scope:            scope,
	}
}

// ParseSubscriptionID parses 'input' into a RoleDefinitionResourceID
func ParseRoleDefinitionResourceId(input string) (*RoleDefinitionResourceId, error) {
	segments := strings.Split(input, "/providers/Microsoft.Authorization/roleDefinitions/")
	switch {
	case strings.HasPrefix(input, "/subscriptions/"):
		return &RoleDefinitionResourceId{
			Scope:            segments[0],
			RoleDefinitionId: segments[1],
		}, nil
	case strings.HasPrefix(input, "/providers/"):
		return &RoleDefinitionResourceId{
			RoleDefinitionId: segments[0],
		}, nil
	default:
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}
}

func (id RoleDefinitionResourceId) ID() string {
	return fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", id.Scope, id.RoleDefinitionId)
}

func (id RoleDefinitionResourceId) String() string {
	components := []string{
		fmt.Sprintf("Role Definition ID: %q", id.RoleDefinitionId),
	}
	if id.Scope != "" {
		components = append(components, fmt.Sprintf("Scope: %q", id.Scope))
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}
