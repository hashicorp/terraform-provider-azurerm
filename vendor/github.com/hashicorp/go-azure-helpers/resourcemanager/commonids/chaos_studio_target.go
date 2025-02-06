// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ChaosStudioTargetId{}

// ChaosStudioTargetId is a struct representing the Resource ID for an App Service Plan
type ChaosStudioTargetId struct {
	Scope      string
	TargetName string
}

// NewChaosStudioTargetID returns a new ChaosStudioTargetId struct
func NewChaosStudioTargetID(scope string, targetName string) ChaosStudioTargetId {
	return ChaosStudioTargetId{
		Scope:      scope,
		TargetName: targetName,
	}
}

// ParseChaosStudioTargetID parses 'input' into a ChaosStudioTargetId
func ParseChaosStudioTargetID(input string) (*ChaosStudioTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChaosStudioTargetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChaosStudioTargetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseChaosStudioTargetIDInsensitively parses 'input' case-insensitively into a ChaosStudioTargetId
// note: this method should only be used for API response data and not user input
func ParseChaosStudioTargetIDInsensitively(input string) (*ChaosStudioTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChaosStudioTargetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChaosStudioTargetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ChaosStudioTargetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.TargetName, ok = input.Parsed["targetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "targetName", input)
	}

	return nil
}

// ValidateChaosStudioTargetID checks that 'input' can be parsed as an App Service Plan ID
func ValidateChaosStudioTargetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseChaosStudioTargetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Service Plan ID
func (id ChaosStudioTargetId) ID() string {
	fmtString := "%s/providers/Microsoft.Chaos/targets/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.TargetName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Service Plan ID
func (id ChaosStudioTargetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticTargets", "targets", "targets"),
		resourceids.UserSpecifiedSegment("targetName", "targetName"),
	}
}

// String returns a human-readable description of this App Service Plan ID
func (id ChaosStudioTargetId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Target Name: %q", id.TargetName),
	}
	return fmt.Sprintf("Chaos Studio Target (%s)", strings.Join(components, "\n"))
}
