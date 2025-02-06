package roleassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleAssignmentId{})
}

var _ resourceids.ResourceId = &ScopedRoleAssignmentId{}

// ScopedRoleAssignmentId is a struct representing the Resource ID for a Scoped Role Assignment
type ScopedRoleAssignmentId struct {
	Scope              string
	RoleAssignmentName string
}

// NewScopedRoleAssignmentID returns a new ScopedRoleAssignmentId struct
func NewScopedRoleAssignmentID(scope string, roleAssignmentName string) ScopedRoleAssignmentId {
	return ScopedRoleAssignmentId{
		Scope:              scope,
		RoleAssignmentName: roleAssignmentName,
	}
}

// ParseScopedRoleAssignmentID parses 'input' into a ScopedRoleAssignmentId
func ParseScopedRoleAssignmentID(input string) (*ScopedRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleAssignmentIDInsensitively parses 'input' case-insensitively into a ScopedRoleAssignmentId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleAssignmentIDInsensitively(input string) (*ScopedRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleAssignmentName, ok = input.Parsed["roleAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentName", input)
	}

	return nil
}

// ValidateScopedRoleAssignmentID checks that 'input' can be parsed as a Scoped Role Assignment ID
func ValidateScopedRoleAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Assignment ID
func (id ScopedRoleAssignmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Assignment ID
func (id ScopedRoleAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleAssignments", "roleAssignments", "roleAssignments"),
		resourceids.UserSpecifiedSegment("roleAssignmentName", "roleAssignmentName"),
	}
}

// String returns a human-readable description of this Scoped Role Assignment ID
func (id ScopedRoleAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Assignment Name: %q", id.RoleAssignmentName),
	}
	return fmt.Sprintf("Scoped Role Assignment (%s)", strings.Join(components, "\n"))
}
