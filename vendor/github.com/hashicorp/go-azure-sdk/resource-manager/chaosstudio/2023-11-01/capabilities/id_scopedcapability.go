package capabilities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedCapabilityId{}

// ScopedCapabilityId is a struct representing the Resource ID for a Scoped Capability
type ScopedCapabilityId struct {
	Scope          string
	TargetName     string
	CapabilityName string
}

// NewScopedCapabilityID returns a new ScopedCapabilityId struct
func NewScopedCapabilityID(scope string, targetName string, capabilityName string) ScopedCapabilityId {
	return ScopedCapabilityId{
		Scope:          scope,
		TargetName:     targetName,
		CapabilityName: capabilityName,
	}
}

// ParseScopedCapabilityID parses 'input' into a ScopedCapabilityId
func ParseScopedCapabilityID(input string) (*ScopedCapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedCapabilityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedCapabilityId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedCapabilityIDInsensitively parses 'input' case-insensitively into a ScopedCapabilityId
// note: this method should only be used for API response data and not user input
func ParseScopedCapabilityIDInsensitively(input string) (*ScopedCapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedCapabilityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedCapabilityId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedCapabilityId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.TargetName, ok = input.Parsed["targetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "targetName", input)
	}

	if id.CapabilityName, ok = input.Parsed["capabilityName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capabilityName", input)
	}

	return nil
}

// ValidateScopedCapabilityID checks that 'input' can be parsed as a Scoped Capability ID
func ValidateScopedCapabilityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedCapabilityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Capability ID
func (id ScopedCapabilityId) ID() string {
	fmtString := "/%s/providers/Microsoft.Chaos/targets/%s/capabilities/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.TargetName, id.CapabilityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Capability ID
func (id ScopedCapabilityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftChaos", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticTargets", "targets", "targets"),
		resourceids.UserSpecifiedSegment("targetName", "targetValue"),
		resourceids.StaticSegment("staticCapabilities", "capabilities", "capabilities"),
		resourceids.UserSpecifiedSegment("capabilityName", "capabilityValue"),
	}
}

// String returns a human-readable description of this Scoped Capability ID
func (id ScopedCapabilityId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Target Name: %q", id.TargetName),
		fmt.Sprintf("Capability Name: %q", id.CapabilityName),
	}
	return fmt.Sprintf("Scoped Capability (%s)", strings.Join(components, "\n"))
}
