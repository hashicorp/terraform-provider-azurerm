package configurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedConfigurationAssignmentId{}

// ScopedConfigurationAssignmentId is a struct representing the Resource ID for a Scoped Configuration Assignment
type ScopedConfigurationAssignmentId struct {
	Scope                       string
	ConfigurationAssignmentName string
}

// NewScopedConfigurationAssignmentID returns a new ScopedConfigurationAssignmentId struct
func NewScopedConfigurationAssignmentID(scope string, configurationAssignmentName string) ScopedConfigurationAssignmentId {
	return ScopedConfigurationAssignmentId{
		Scope:                       scope,
		ConfigurationAssignmentName: configurationAssignmentName,
	}
}

// ParseScopedConfigurationAssignmentID parses 'input' into a ScopedConfigurationAssignmentId
func ParseScopedConfigurationAssignmentID(input string) (*ScopedConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedConfigurationAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ConfigurationAssignmentName, ok = parsed.Parsed["configurationAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationAssignmentName", *parsed)
	}

	return &id, nil
}

// ParseScopedConfigurationAssignmentIDInsensitively parses 'input' case-insensitively into a ScopedConfigurationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseScopedConfigurationAssignmentIDInsensitively(input string) (*ScopedConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedConfigurationAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ConfigurationAssignmentName, ok = parsed.Parsed["configurationAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationAssignmentName", *parsed)
	}

	return &id, nil
}

// ValidateScopedConfigurationAssignmentID checks that 'input' can be parsed as a Scoped Configuration Assignment ID
func ValidateScopedConfigurationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedConfigurationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Configuration Assignment ID
func (id ScopedConfigurationAssignmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Maintenance/configurationAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.ConfigurationAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Configuration Assignment ID
func (id ScopedConfigurationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticConfigurationAssignments", "configurationAssignments", "configurationAssignments"),
		resourceids.UserSpecifiedSegment("configurationAssignmentName", "configurationAssignmentValue"),
	}
}

// String returns a human-readable description of this Scoped Configuration Assignment ID
func (id ScopedConfigurationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Configuration Assignment Name: %q", id.ConfigurationAssignmentName),
	}
	return fmt.Sprintf("Scoped Configuration Assignment (%s)", strings.Join(components, "\n"))
}
