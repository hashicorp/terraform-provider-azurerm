package roledefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedRoleDefinitionId{}

// ScopedRoleDefinitionId is a struct representing the Resource ID for a Scoped Role Definition
type ScopedRoleDefinitionId struct {
	Scope            string
	RoleDefinitionId string
}

// NewScopedRoleDefinitionID returns a new ScopedRoleDefinitionId struct
func NewScopedRoleDefinitionID(scope string, roleDefinitionId string) ScopedRoleDefinitionId {
	return ScopedRoleDefinitionId{
		Scope:            scope,
		RoleDefinitionId: roleDefinitionId,
	}
}

// ParseScopedRoleDefinitionID parses 'input' into a ScopedRoleDefinitionId
func ParseScopedRoleDefinitionID(input string) (*ScopedRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleDefinitionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleDefinitionId, ok = parsed.Parsed["roleDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionId", *parsed)
	}

	return &id, nil
}

// ParseScopedRoleDefinitionIDInsensitively parses 'input' case-insensitively into a ScopedRoleDefinitionId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleDefinitionIDInsensitively(input string) (*ScopedRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleDefinitionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleDefinitionId, ok = parsed.Parsed["roleDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionId", *parsed)
	}

	return &id, nil
}

// ValidateScopedRoleDefinitionID checks that 'input' can be parsed as a Scoped Role Definition ID
func ValidateScopedRoleDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Definition ID
func (id ScopedRoleDefinitionId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleDefinitions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleDefinitionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Definition ID
func (id ScopedRoleDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleDefinitions", "roleDefinitions", "roleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionId", "roleDefinitionIdValue"),
	}
}

// String returns a human-readable description of this Scoped Role Definition ID
func (id ScopedRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Definition: %q", id.RoleDefinitionId),
	}
	return fmt.Sprintf("Scoped Role Definition (%s)", strings.Join(components, "\n"))
}
