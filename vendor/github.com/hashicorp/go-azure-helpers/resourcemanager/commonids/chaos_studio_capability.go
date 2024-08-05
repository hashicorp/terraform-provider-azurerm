// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ChaosStudioCapabilityId{}

// ChaosStudioCapabilityId is a struct representing the Resource ID for an App Service Plan
type ChaosStudioCapabilityId struct {
	Scope          string
	TargetName     string
	CapabilityName string
}

// NewChaosStudioCapabilityID returns a new ChaosStudioCapabilityId struct
func NewChaosStudioCapabilityID(scope string, targetName string, capabilityName string) ChaosStudioCapabilityId {
	return ChaosStudioCapabilityId{
		Scope:          scope,
		TargetName:     targetName,
		CapabilityName: capabilityName,
	}
}

// ParseChaosStudioCapabilityID parses 'input' into a ChaosStudioCapabilityId
func ParseChaosStudioCapabilityID(input string) (*ChaosStudioCapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChaosStudioCapabilityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChaosStudioCapabilityId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseChaosStudioCapabilityIDInsensitively parses 'input' case-insensitively into a ChaosStudioCapabilityId
// note: this method should only be used for API response data and not user input
func ParseChaosStudioCapabilityIDInsensitively(input string) (*ChaosStudioCapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChaosStudioCapabilityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChaosStudioCapabilityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ChaosStudioCapabilityId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateChaosStudioCapabilityID checks that 'input' can be parsed as an App Service Plan ID
func ValidateChaosStudioCapabilityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseChaosStudioCapabilityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Service Plan ID
func (id ChaosStudioCapabilityId) ID() string {
	fmtString := "%s/providers/Microsoft.Chaos/targets/%s/capabilities/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.TargetName, id.CapabilityName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Service Plan ID
func (id ChaosStudioCapabilityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticTargets", "targets", "targets"),
		resourceids.UserSpecifiedSegment("targetName", "targetName"),
		resourceids.StaticSegment("staticCapabilities", "capabilities", "capabilities"),
		resourceids.UserSpecifiedSegment("capabilityName", "capabilityName"),
	}
}

// String returns a human-readable description of this App Service Plan ID
func (id ChaosStudioCapabilityId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Target Name: %q", id.TargetName),
		fmt.Sprintf("Capability Name: %q", id.TargetName),
	}
	return fmt.Sprintf("Chaos Studio Capability (%s)", strings.Join(components, "\n"))
}
