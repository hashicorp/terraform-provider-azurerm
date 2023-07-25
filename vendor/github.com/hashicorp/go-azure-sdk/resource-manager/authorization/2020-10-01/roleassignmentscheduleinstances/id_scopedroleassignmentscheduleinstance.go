package roleassignmentscheduleinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedRoleAssignmentScheduleInstanceId{}

// ScopedRoleAssignmentScheduleInstanceId is a struct representing the Resource ID for a Scoped Role Assignment Schedule Instance
type ScopedRoleAssignmentScheduleInstanceId struct {
	Scope                              string
	RoleAssignmentScheduleInstanceName string
}

// NewScopedRoleAssignmentScheduleInstanceID returns a new ScopedRoleAssignmentScheduleInstanceId struct
func NewScopedRoleAssignmentScheduleInstanceID(scope string, roleAssignmentScheduleInstanceName string) ScopedRoleAssignmentScheduleInstanceId {
	return ScopedRoleAssignmentScheduleInstanceId{
		Scope:                              scope,
		RoleAssignmentScheduleInstanceName: roleAssignmentScheduleInstanceName,
	}
}

// ParseScopedRoleAssignmentScheduleInstanceID parses 'input' into a ScopedRoleAssignmentScheduleInstanceId
func ParseScopedRoleAssignmentScheduleInstanceID(input string) (*ScopedRoleAssignmentScheduleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleAssignmentScheduleInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleAssignmentScheduleInstanceId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleAssignmentScheduleInstanceName, ok = parsed.Parsed["roleAssignmentScheduleInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentScheduleInstanceName", *parsed)
	}

	return &id, nil
}

// ParseScopedRoleAssignmentScheduleInstanceIDInsensitively parses 'input' case-insensitively into a ScopedRoleAssignmentScheduleInstanceId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleAssignmentScheduleInstanceIDInsensitively(input string) (*ScopedRoleAssignmentScheduleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleAssignmentScheduleInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleAssignmentScheduleInstanceId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleAssignmentScheduleInstanceName, ok = parsed.Parsed["roleAssignmentScheduleInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentScheduleInstanceName", *parsed)
	}

	return &id, nil
}

// ValidateScopedRoleAssignmentScheduleInstanceID checks that 'input' can be parsed as a Scoped Role Assignment Schedule Instance ID
func ValidateScopedRoleAssignmentScheduleInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleAssignmentScheduleInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Assignment Schedule Instance ID
func (id ScopedRoleAssignmentScheduleInstanceId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleAssignmentScheduleInstances/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleAssignmentScheduleInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Assignment Schedule Instance ID
func (id ScopedRoleAssignmentScheduleInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleAssignmentScheduleInstances", "roleAssignmentScheduleInstances", "roleAssignmentScheduleInstances"),
		resourceids.UserSpecifiedSegment("roleAssignmentScheduleInstanceName", "roleAssignmentScheduleInstanceValue"),
	}
}

// String returns a human-readable description of this Scoped Role Assignment Schedule Instance ID
func (id ScopedRoleAssignmentScheduleInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Assignment Schedule Instance Name: %q", id.RoleAssignmentScheduleInstanceName),
	}
	return fmt.Sprintf("Scoped Role Assignment Schedule Instance (%s)", strings.Join(components, "\n"))
}
