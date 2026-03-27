// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

var roleDefinitionIdPattern = regexp.MustCompile(`(?i)^(.*)/providers/Microsoft\.Authorization/roleDefinitions/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`)

type PimRoleDefinitionId struct {
	Scope              string
	RoleDefinitionUuid string
}

func NewPimRoleDefinitionID(scope string, roleDefinitionUuid string) PimRoleDefinitionId {
	return PimRoleDefinitionId{
		Scope:              scope,
		RoleDefinitionUuid: roleDefinitionUuid,
	}
}

func (id PimRoleDefinitionId) ID() string {
	if id.Scope == "" {
		return fmt.Sprintf("/providers/Microsoft.Authorization/roleDefinitions/%s", id.RoleDefinitionUuid)
	}
	return fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", id.Scope, id.RoleDefinitionUuid)
}

func (id PimRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Definition: %q", id.RoleDefinitionUuid),
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}

// PimRoleDefinitionID parses a role definition ID that may be either scoped
// (e.g., "/subscriptions/{sub}/providers/Microsoft.Authorization/roleDefinitions/{uuid}")
// or unscoped (e.g., "/providers/Microsoft.Authorization/roleDefinitions/{uuid}").
func PimRoleDefinitionID(input string) (*PimRoleDefinitionId, error) {
	matches := roleDefinitionIdPattern.FindStringSubmatch(input)
	if len(matches) != 3 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}

	return &PimRoleDefinitionId{
		Scope:              matches[1],
		RoleDefinitionUuid: strings.ToLower(matches[2]),
	}, nil
}

// PimRoleDefinitionIdsMatch compares two role definition IDs by their UUIDs,
// ignoring differences in scope prefix.
func PimRoleDefinitionIdsMatch(a, b string) bool {
	idA, errA := PimRoleDefinitionID(a)
	idB, errB := PimRoleDefinitionID(b)
	if errA != nil || errB != nil {
		return strings.EqualFold(a, b)
	}

	return strings.EqualFold(idA.RoleDefinitionUuid, idB.RoleDefinitionUuid)
}
