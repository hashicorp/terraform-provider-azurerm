// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

var roleDefinitionResourceIdPattern = regexp.MustCompile(`(?i)^(.*)/providers/Microsoft\.Authorization/roleDefinitions/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`)

type RoleDefinitionResourceId struct {
	Scope              string
	RoleDefinitionUuid string
}

func NewRoleDefinitionResourceID(scope string, roleDefinitionUuid string) RoleDefinitionResourceId {
	return RoleDefinitionResourceId{
		Scope:              scope,
		RoleDefinitionUuid: roleDefinitionUuid,
	}
}

func (id RoleDefinitionResourceId) ID() string {
	if id.Scope == "" {
		return fmt.Sprintf("/providers/Microsoft.Authorization/roleDefinitions/%s", id.RoleDefinitionUuid)
	}
	return fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", id.Scope, id.RoleDefinitionUuid)
}

func (id RoleDefinitionResourceId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Definition: %q", id.RoleDefinitionUuid),
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}

// RoleDefinitionResourceID parses a role definition ID with a scope prefix such as a subscription,
// resource group, resource, or the provider-level /providers prefix.
// The SDK's roledefinitions.ParseScopedRoleDefinitionID parser is not used here because it does not accept
// provider-level role definition IDs that start with /providers/Microsoft.Authorization.
func RoleDefinitionResourceID(input string) (*RoleDefinitionResourceId, error) {
	matches := roleDefinitionResourceIdPattern.FindStringSubmatch(input)
	// FindStringSubmatch returns 3 values: the full ID match, then the scope and UUID capture groups.
	if len(matches) != 3 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}

	return &RoleDefinitionResourceId{
		Scope:              matches[1],
		RoleDefinitionUuid: strings.ToLower(matches[2]),
	}, nil
}

// RoleDefinitionResourceIdsMatch compares two role definition IDs by their UUIDs,
// ignoring differences in scope prefix.
func RoleDefinitionResourceIdsMatch(a, b string) bool {
	idA, errA := RoleDefinitionResourceID(a)
	idB, errB := RoleDefinitionResourceID(b)
	if errA != nil || errB != nil {
		return strings.EqualFold(a, b)
	}

	return strings.EqualFold(idA.RoleDefinitionUuid, idB.RoleDefinitionUuid)
}
