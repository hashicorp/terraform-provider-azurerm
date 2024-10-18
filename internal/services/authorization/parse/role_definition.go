// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RoleDefinitionID struct {
	ResourceID string
	Scope      string
	RoleID     string
}

var _ resourceids.Id = RoleDefinitionID{}

func (r RoleDefinitionID) ID() string {
	return fmt.Sprintf("%s|%s", r.ResourceID, r.Scope)
}

func (r RoleDefinitionID) String() string {
	components := []string{
		fmt.Sprintf("Resource ID: %q", r.ResourceID),
		fmt.Sprintf("Scope: %q", r.Scope),
		fmt.Sprintf("Role Definition: %q", r.RoleID),
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}

func (r RoleDefinitionID) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleDefinitions", "roleDefinitions", "roleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionId", "roleDefinitionIdValue"),
	}
}

var roleDefinitionIdRegexp *regexp.Regexp = regexp.MustCompile(
	"^(?i)(?:/subscriptions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?" +
		"/providers/Microsoft\\.Authorization/roleDefinitions/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$")

var subscriptionScopeRegexp *regexp.Regexp = regexp.MustCompile(
	"^(?i)(?:/subscriptions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}" +
		"(?:/resourceGroups/[_\\-.()0-9a-zA-Z]{1,89}[_\\-()0-9a-zA-Z]{1}(?:/providers/.+)?)?)$")

var mgmtGroupScopeRegexp *regexp.Regexp = regexp.MustCompile(
	"^(?i)(?:/providers/Microsoft\\.Management/managementGroups/[_\\-.()0-9a-zA-Z]{1,89}[_\\-()0-9a-zA-Z]{1}" +
		"(?:/subscriptions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}" +
		"(?:/resourceGroups/[_\\-.()0-9a-zA-Z]{1,89}[_\\-()0-9a-zA-Z]{1}(?:/providers/.+)?)?)?)$")

const allScopesToken string = "/"

// RoleDefinitionId is a pseudo ID for storing Scope parameter as this it not retrievable from API
// It is formed of the Azure Resource ID for the Role and the Scope it is created against
func RoleDefinitionId(input string) (*RoleDefinitionID, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}

	roleDefinition := roleDefinitionIdRegexp.FindStringSubmatch(parts[0])
	if len(roleDefinition) != 2 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", parts[0])
	}

	var scope string = subscriptionScopeRegexp.FindString(parts[1])
	if len(scope) == 0 {
		scope = mgmtGroupScopeRegexp.FindString(parts[1])
		if len(scope) == 0 {
			scope = parts[1]
			if scope != allScopesToken {
				return nil, fmt.Errorf("could not parse scope from Role Definition ID, invalid format %q", parts[1])
			}
		}
	}

	roleDefinitionID := RoleDefinitionID{
		ResourceID: roleDefinition[0],
		Scope:      scope,
		RoleID:     roleDefinition[1],
	}

	return &roleDefinitionID, nil
}
