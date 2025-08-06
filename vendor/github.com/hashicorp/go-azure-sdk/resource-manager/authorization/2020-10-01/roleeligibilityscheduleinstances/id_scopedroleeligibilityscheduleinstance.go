package roleeligibilityscheduleinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleEligibilityScheduleInstanceId{})
}

var _ resourceids.ResourceId = &ScopedRoleEligibilityScheduleInstanceId{}

// ScopedRoleEligibilityScheduleInstanceId is a struct representing the Resource ID for a Scoped Role Eligibility Schedule Instance
type ScopedRoleEligibilityScheduleInstanceId struct {
	Scope                               string
	RoleEligibilityScheduleInstanceName string
}

// NewScopedRoleEligibilityScheduleInstanceID returns a new ScopedRoleEligibilityScheduleInstanceId struct
func NewScopedRoleEligibilityScheduleInstanceID(scope string, roleEligibilityScheduleInstanceName string) ScopedRoleEligibilityScheduleInstanceId {
	return ScopedRoleEligibilityScheduleInstanceId{
		Scope:                               scope,
		RoleEligibilityScheduleInstanceName: roleEligibilityScheduleInstanceName,
	}
}

// ParseScopedRoleEligibilityScheduleInstanceID parses 'input' into a ScopedRoleEligibilityScheduleInstanceId
func ParseScopedRoleEligibilityScheduleInstanceID(input string) (*ScopedRoleEligibilityScheduleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleEligibilityScheduleInstanceIDInsensitively parses 'input' case-insensitively into a ScopedRoleEligibilityScheduleInstanceId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleEligibilityScheduleInstanceIDInsensitively(input string) (*ScopedRoleEligibilityScheduleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleEligibilityScheduleInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleEligibilityScheduleInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleEligibilityScheduleInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleEligibilityScheduleInstanceName, ok = input.Parsed["roleEligibilityScheduleInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleEligibilityScheduleInstanceName", input)
	}

	return nil
}

// ValidateScopedRoleEligibilityScheduleInstanceID checks that 'input' can be parsed as a Scoped Role Eligibility Schedule Instance ID
func ValidateScopedRoleEligibilityScheduleInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleEligibilityScheduleInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Eligibility Schedule Instance ID
func (id ScopedRoleEligibilityScheduleInstanceId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleEligibilityScheduleInstances/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleEligibilityScheduleInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Eligibility Schedule Instance ID
func (id ScopedRoleEligibilityScheduleInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleEligibilityScheduleInstances", "roleEligibilityScheduleInstances", "roleEligibilityScheduleInstances"),
		resourceids.UserSpecifiedSegment("roleEligibilityScheduleInstanceName", "roleEligibilityScheduleInstanceName"),
	}
}

// String returns a human-readable description of this Scoped Role Eligibility Schedule Instance ID
func (id ScopedRoleEligibilityScheduleInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Eligibility Schedule Instance Name: %q", id.RoleEligibilityScheduleInstanceName),
	}
	return fmt.Sprintf("Scoped Role Eligibility Schedule Instance (%s)", strings.Join(components, "\n"))
}
