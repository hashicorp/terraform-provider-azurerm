// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
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

// RoleDefinitionId is a pseudo ID for storing Scope parameter as this it not retrievable from API
// It is formed of the Azure Resource ID for the Role and the Scope it is created against
func RoleDefinitionId(input string) (*RoleDefinitionID, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}

	idParts := strings.Split(parts[0], "roleDefinitions/")

	if !strings.HasPrefix(parts[1], "/subscriptions/") && !strings.HasPrefix(parts[1], "/providers/Microsoft.Management/managementGroups/") {
		return nil, fmt.Errorf("failed to parse scope from Role Definition ID %q", input)
	}

	roleDefinitionID := RoleDefinitionID{
		ResourceID: parts[0],
		Scope:      parts[1],
	}

	if len(idParts) < 1 {
		return nil, fmt.Errorf("failed to parse Role Definition ID from resource ID %q", input)
	} else {
		roleDefinitionID.RoleID = idParts[1]
	}

	return &roleDefinitionID, nil
}
