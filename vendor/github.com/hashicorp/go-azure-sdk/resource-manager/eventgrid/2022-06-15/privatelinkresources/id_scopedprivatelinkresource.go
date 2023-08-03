package privatelinkresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedPrivateLinkResourceId{}

// ScopedPrivateLinkResourceId is a struct representing the Resource ID for a Scoped Private Link Resource
type ScopedPrivateLinkResourceId struct {
	Scope                   string
	PrivateLinkResourceName string
}

// NewScopedPrivateLinkResourceID returns a new ScopedPrivateLinkResourceId struct
func NewScopedPrivateLinkResourceID(scope string, privateLinkResourceName string) ScopedPrivateLinkResourceId {
	return ScopedPrivateLinkResourceId{
		Scope:                   scope,
		PrivateLinkResourceName: privateLinkResourceName,
	}
}

// ParseScopedPrivateLinkResourceID parses 'input' into a ScopedPrivateLinkResourceId
func ParseScopedPrivateLinkResourceID(input string) (*ScopedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedPrivateLinkResourceId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.PrivateLinkResourceName, ok = parsed.Parsed["privateLinkResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkResourceName", *parsed)
	}

	return &id, nil
}

// ParseScopedPrivateLinkResourceIDInsensitively parses 'input' case-insensitively into a ScopedPrivateLinkResourceId
// note: this method should only be used for API response data and not user input
func ParseScopedPrivateLinkResourceIDInsensitively(input string) (*ScopedPrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedPrivateLinkResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedPrivateLinkResourceId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.PrivateLinkResourceName, ok = parsed.Parsed["privateLinkResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkResourceName", *parsed)
	}

	return &id, nil
}

// ValidateScopedPrivateLinkResourceID checks that 'input' can be parsed as a Scoped Private Link Resource ID
func ValidateScopedPrivateLinkResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedPrivateLinkResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Private Link Resource ID
func (id ScopedPrivateLinkResourceId) ID() string {
	fmtString := "/%s/privateLinkResources/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.PrivateLinkResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Private Link Resource ID
func (id ScopedPrivateLinkResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticPrivateLinkResources", "privateLinkResources", "privateLinkResources"),
		resourceids.UserSpecifiedSegment("privateLinkResourceName", "privateLinkResourceValue"),
	}
}

// String returns a human-readable description of this Scoped Private Link Resource ID
func (id ScopedPrivateLinkResourceId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Private Link Resource Name: %q", id.PrivateLinkResourceName),
	}
	return fmt.Sprintf("Scoped Private Link Resource (%s)", strings.Join(components, "\n"))
}
