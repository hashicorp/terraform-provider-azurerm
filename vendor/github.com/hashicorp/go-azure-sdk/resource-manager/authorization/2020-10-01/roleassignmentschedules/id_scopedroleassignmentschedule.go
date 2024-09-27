package roleassignmentschedules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleAssignmentScheduleId{})
}

var _ resourceids.ResourceId = &ScopedRoleAssignmentScheduleId{}

// ScopedRoleAssignmentScheduleId is a struct representing the Resource ID for a Scoped Role Assignment Schedule
type ScopedRoleAssignmentScheduleId struct {
	Scope                      string
	RoleAssignmentScheduleName string
}

// NewScopedRoleAssignmentScheduleID returns a new ScopedRoleAssignmentScheduleId struct
func NewScopedRoleAssignmentScheduleID(scope string, roleAssignmentScheduleName string) ScopedRoleAssignmentScheduleId {
	return ScopedRoleAssignmentScheduleId{
		Scope:                      scope,
		RoleAssignmentScheduleName: roleAssignmentScheduleName,
	}
}

// ParseScopedRoleAssignmentScheduleID parses 'input' into a ScopedRoleAssignmentScheduleId
func ParseScopedRoleAssignmentScheduleID(input string) (*ScopedRoleAssignmentScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleAssignmentScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleAssignmentScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleAssignmentScheduleIDInsensitively parses 'input' case-insensitively into a ScopedRoleAssignmentScheduleId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleAssignmentScheduleIDInsensitively(input string) (*ScopedRoleAssignmentScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleAssignmentScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleAssignmentScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleAssignmentScheduleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleAssignmentScheduleName, ok = input.Parsed["roleAssignmentScheduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentScheduleName", input)
	}

	return nil
}

// ValidateScopedRoleAssignmentScheduleID checks that 'input' can be parsed as a Scoped Role Assignment Schedule ID
func ValidateScopedRoleAssignmentScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleAssignmentScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Assignment Schedule ID
func (id ScopedRoleAssignmentScheduleId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleAssignmentSchedules/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleAssignmentScheduleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Assignment Schedule ID
func (id ScopedRoleAssignmentScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleAssignmentSchedules", "roleAssignmentSchedules", "roleAssignmentSchedules"),
		resourceids.UserSpecifiedSegment("roleAssignmentScheduleName", "roleAssignmentScheduleName"),
	}
}

// String returns a human-readable description of this Scoped Role Assignment Schedule ID
func (id ScopedRoleAssignmentScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Assignment Schedule Name: %q", id.RoleAssignmentScheduleName),
	}
	return fmt.Sprintf("Scoped Role Assignment Schedule (%s)", strings.Join(components, "\n"))
}
