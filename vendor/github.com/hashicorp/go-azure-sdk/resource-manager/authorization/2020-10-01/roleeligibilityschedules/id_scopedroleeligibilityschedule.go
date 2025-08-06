package roleeligibilityschedules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleEligibilityScheduleId{})
}

var _ resourceids.ResourceId = &ScopedRoleEligibilityScheduleId{}

// ScopedRoleEligibilityScheduleId is a struct representing the Resource ID for a Scoped Role Eligibility Schedule
type ScopedRoleEligibilityScheduleId struct {
	Scope                       string
	RoleEligibilityScheduleName string
}

// NewScopedRoleEligibilityScheduleID returns a new ScopedRoleEligibilityScheduleId struct
func NewScopedRoleEligibilityScheduleID(scope string, roleEligibilityScheduleName string) ScopedRoleEligibilityScheduleId {
	return ScopedRoleEligibilityScheduleId{
		Scope:                       scope,
		RoleEligibilityScheduleName: roleEligibilityScheduleName,
	}
}

// ParseScopedRoleEligibilityScheduleID parses 'input' into a ScopedRoleEligibilityScheduleId
func ParseScopedRoleEligibilityScheduleID(input string) (*ScopedRoleEligibilityScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleEligibilityScheduleIDInsensitively parses 'input' case-insensitively into a ScopedRoleEligibilityScheduleId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleEligibilityScheduleIDInsensitively(input string) (*ScopedRoleEligibilityScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleEligibilityScheduleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleEligibilityScheduleName, ok = input.Parsed["roleEligibilityScheduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleEligibilityScheduleName", input)
	}

	return nil
}

// ValidateScopedRoleEligibilityScheduleID checks that 'input' can be parsed as a Scoped Role Eligibility Schedule ID
func ValidateScopedRoleEligibilityScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleEligibilityScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Eligibility Schedule ID
func (id ScopedRoleEligibilityScheduleId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleEligibilitySchedules/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleEligibilityScheduleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Eligibility Schedule ID
func (id ScopedRoleEligibilityScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleEligibilitySchedules", "roleEligibilitySchedules", "roleEligibilitySchedules"),
		resourceids.UserSpecifiedSegment("roleEligibilityScheduleName", "roleEligibilityScheduleName"),
	}
}

// String returns a human-readable description of this Scoped Role Eligibility Schedule ID
func (id ScopedRoleEligibilityScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Eligibility Schedule Name: %q", id.RoleEligibilityScheduleName),
	}
	return fmt.Sprintf("Scoped Role Eligibility Schedule (%s)", strings.Join(components, "\n"))
}
