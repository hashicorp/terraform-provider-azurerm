package extensions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedExtensionId{}

// ScopedExtensionId is a struct representing the Resource ID for a Scoped Extension
type ScopedExtensionId struct {
	Scope         string
	ExtensionName string
}

// NewScopedExtensionID returns a new ScopedExtensionId struct
func NewScopedExtensionID(scope string, extensionName string) ScopedExtensionId {
	return ScopedExtensionId{
		Scope:         scope,
		ExtensionName: extensionName,
	}
}

// ParseScopedExtensionID parses 'input' into a ScopedExtensionId
func ParseScopedExtensionID(input string) (*ScopedExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedExtensionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "extensionName", *parsed)
	}

	return &id, nil
}

// ParseScopedExtensionIDInsensitively parses 'input' case-insensitively into a ScopedExtensionId
// note: this method should only be used for API response data and not user input
func ParseScopedExtensionIDInsensitively(input string) (*ScopedExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedExtensionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "extensionName", *parsed)
	}

	return &id, nil
}

// ValidateScopedExtensionID checks that 'input' can be parsed as a Scoped Extension ID
func ValidateScopedExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Extension ID
func (id ScopedExtensionId) ID() string {
	fmtString := "/%s/providers/Microsoft.KubernetesConfiguration/extensions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.ExtensionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Extension ID
func (id ScopedExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetesConfiguration", "Microsoft.KubernetesConfiguration", "Microsoft.KubernetesConfiguration"),
		resourceids.StaticSegment("staticExtensions", "extensions", "extensions"),
		resourceids.UserSpecifiedSegment("extensionName", "extensionValue"),
	}
}

// String returns a human-readable description of this Scoped Extension ID
func (id ScopedExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Extension Name: %q", id.ExtensionName),
	}
	return fmt.Sprintf("Scoped Extension (%s)", strings.Join(components, "\n"))
}
