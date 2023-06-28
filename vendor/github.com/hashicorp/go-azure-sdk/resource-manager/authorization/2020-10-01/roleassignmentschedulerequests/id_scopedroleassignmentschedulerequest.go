package roleassignmentschedulerequests

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedRoleAssignmentScheduleRequestId{}

// ScopedRoleAssignmentScheduleRequestId is a struct representing the Resource ID for a Scoped Role Assignment Schedule Request
type ScopedRoleAssignmentScheduleRequestId struct {
	Scope                             string
	RoleAssignmentScheduleRequestName string
}

// NewScopedRoleAssignmentScheduleRequestID returns a new ScopedRoleAssignmentScheduleRequestId struct
func NewScopedRoleAssignmentScheduleRequestID(scope string, roleAssignmentScheduleRequestName string) ScopedRoleAssignmentScheduleRequestId {
	return ScopedRoleAssignmentScheduleRequestId{
		Scope:                             scope,
		RoleAssignmentScheduleRequestName: roleAssignmentScheduleRequestName,
	}
}

// ParseScopedRoleAssignmentScheduleRequestID parses 'input' into a ScopedRoleAssignmentScheduleRequestId
func ParseScopedRoleAssignmentScheduleRequestID(input string) (*ScopedRoleAssignmentScheduleRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleAssignmentScheduleRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleAssignmentScheduleRequestId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleAssignmentScheduleRequestName, ok = parsed.Parsed["roleAssignmentScheduleRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentScheduleRequestName", *parsed)
	}

	return &id, nil
}

// ParseScopedRoleAssignmentScheduleRequestIDInsensitively parses 'input' case-insensitively into a ScopedRoleAssignmentScheduleRequestId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleAssignmentScheduleRequestIDInsensitively(input string) (*ScopedRoleAssignmentScheduleRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRoleAssignmentScheduleRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRoleAssignmentScheduleRequestId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RoleAssignmentScheduleRequestName, ok = parsed.Parsed["roleAssignmentScheduleRequestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentScheduleRequestName", *parsed)
	}

	return &id, nil
}

// ValidateScopedRoleAssignmentScheduleRequestID checks that 'input' can be parsed as a Scoped Role Assignment Schedule Request ID
func ValidateScopedRoleAssignmentScheduleRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleAssignmentScheduleRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Assignment Schedule Request ID
func (id ScopedRoleAssignmentScheduleRequestId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleAssignmentScheduleRequests/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleAssignmentScheduleRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Assignment Schedule Request ID
func (id ScopedRoleAssignmentScheduleRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleAssignmentScheduleRequests", "roleAssignmentScheduleRequests", "roleAssignmentScheduleRequests"),
		resourceids.UserSpecifiedSegment("roleAssignmentScheduleRequestName", "roleAssignmentScheduleRequestValue"),
	}
}

// String returns a human-readable description of this Scoped Role Assignment Schedule Request ID
func (id ScopedRoleAssignmentScheduleRequestId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Assignment Schedule Request Name: %q", id.RoleAssignmentScheduleRequestName),
	}
	return fmt.Sprintf("Scoped Role Assignment Schedule Request (%s)", strings.Join(components, "\n"))
}
