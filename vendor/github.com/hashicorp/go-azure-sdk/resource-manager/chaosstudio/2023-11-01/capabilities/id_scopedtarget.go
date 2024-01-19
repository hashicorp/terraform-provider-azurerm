package capabilities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedTargetId{}

// ScopedTargetId is a struct representing the Resource ID for a Scoped Target
type ScopedTargetId struct {
	Scope      string
	TargetName string
}

// NewScopedTargetID returns a new ScopedTargetId struct
func NewScopedTargetID(scope string, targetName string) ScopedTargetId {
	return ScopedTargetId{
		Scope:      scope,
		TargetName: targetName,
	}
}

// ParseScopedTargetID parses 'input' into a ScopedTargetId
func ParseScopedTargetID(input string) (*ScopedTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedTargetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedTargetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedTargetIDInsensitively parses 'input' case-insensitively into a ScopedTargetId
// note: this method should only be used for API response data and not user input
func ParseScopedTargetIDInsensitively(input string) (*ScopedTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedTargetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedTargetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedTargetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.TargetName, ok = input.Parsed["targetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "targetName", input)
	}

	return nil
}

// ValidateScopedTargetID checks that 'input' can be parsed as a Scoped Target ID
func ValidateScopedTargetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedTargetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Target ID
func (id ScopedTargetId) ID() string {
	fmtString := "/%s/providers/Microsoft.Chaos/targets/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.TargetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Target ID
func (id ScopedTargetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftChaos", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticTargets", "targets", "targets"),
		resourceids.UserSpecifiedSegment("targetName", "targetValue"),
	}
}

// String returns a human-readable description of this Scoped Target ID
func (id ScopedTargetId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Target Name: %q", id.TargetName),
	}
	return fmt.Sprintf("Scoped Target (%s)", strings.Join(components, "\n"))
}
