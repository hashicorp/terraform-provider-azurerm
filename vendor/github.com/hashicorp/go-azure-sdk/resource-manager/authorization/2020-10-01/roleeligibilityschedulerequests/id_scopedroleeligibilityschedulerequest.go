package roleeligibilityschedulerequests

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleEligibilityScheduleRequestId{})
}

var _ resourceids.ResourceId = &ScopedRoleEligibilityScheduleRequestId{}

// ScopedRoleEligibilityScheduleRequestId is a struct representing the Resource ID for a Scoped Role Eligibility Schedule Request
type ScopedRoleEligibilityScheduleRequestId struct {
	Scope                              string
	RoleEligibilityScheduleRequestName string
}

// NewScopedRoleEligibilityScheduleRequestID returns a new ScopedRoleEligibilityScheduleRequestId struct
func NewScopedRoleEligibilityScheduleRequestID(scope string, roleEligibilityScheduleRequestName string) ScopedRoleEligibilityScheduleRequestId {
	return ScopedRoleEligibilityScheduleRequestId{
		Scope:                              scope,
		RoleEligibilityScheduleRequestName: roleEligibilityScheduleRequestName,
	}
}

// ParseScopedRoleEligibilityScheduleRequestID parses 'input' into a ScopedRoleEligibilityScheduleRequestId
func ParseScopedRoleEligibilityScheduleRequestID(input string) (*ScopedRoleEligibilityScheduleRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleEligibilityScheduleRequestIDInsensitively parses 'input' case-insensitively into a ScopedRoleEligibilityScheduleRequestId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleEligibilityScheduleRequestIDInsensitively(input string) (*ScopedRoleEligibilityScheduleRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleRequestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleEligibilityScheduleRequestId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleEligibilityScheduleRequestName, ok = input.Parsed["roleEligibilityScheduleRequestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleEligibilityScheduleRequestName", input)
	}

	return nil
}

// ValidateScopedRoleEligibilityScheduleRequestID checks that 'input' can be parsed as a Scoped Role Eligibility Schedule Request ID
func ValidateScopedRoleEligibilityScheduleRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleEligibilityScheduleRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Eligibility Schedule Request ID
func (id ScopedRoleEligibilityScheduleRequestId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleEligibilityScheduleRequests/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleEligibilityScheduleRequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Eligibility Schedule Request ID
func (id ScopedRoleEligibilityScheduleRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleEligibilityScheduleRequests", "roleEligibilityScheduleRequests", "roleEligibilityScheduleRequests"),
		resourceids.UserSpecifiedSegment("roleEligibilityScheduleRequestName", "roleEligibilityScheduleRequestName"),
	}
}

// String returns a human-readable description of this Scoped Role Eligibility Schedule Request ID
func (id ScopedRoleEligibilityScheduleRequestId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Eligibility Schedule Request Name: %q", id.RoleEligibilityScheduleRequestName),
	}
	return fmt.Sprintf("Scoped Role Eligibility Schedule Request (%s)", strings.Join(components, "\n"))
}
