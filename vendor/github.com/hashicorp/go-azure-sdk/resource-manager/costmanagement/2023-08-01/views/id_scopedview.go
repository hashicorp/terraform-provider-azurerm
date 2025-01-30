package views

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedViewId{})
}

var _ resourceids.ResourceId = &ScopedViewId{}

// ScopedViewId is a struct representing the Resource ID for a Scoped View
type ScopedViewId struct {
	Scope    string
	ViewName string
}

// NewScopedViewID returns a new ScopedViewId struct
func NewScopedViewID(scope string, viewName string) ScopedViewId {
	return ScopedViewId{
		Scope:    scope,
		ViewName: viewName,
	}
}

// ParseScopedViewID parses 'input' into a ScopedViewId
func ParseScopedViewID(input string) (*ScopedViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedViewId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedViewIDInsensitively parses 'input' case-insensitively into a ScopedViewId
// note: this method should only be used for API response data and not user input
func ParseScopedViewIDInsensitively(input string) (*ScopedViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedViewId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedViewId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.ViewName, ok = input.Parsed["viewName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "viewName", input)
	}

	return nil
}

// ValidateScopedViewID checks that 'input' can be parsed as a Scoped View ID
func ValidateScopedViewID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedViewID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped View ID
func (id ScopedViewId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.ViewName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped View ID
func (id ScopedViewId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticViews", "views", "views"),
		resourceids.UserSpecifiedSegment("viewName", "viewName"),
	}
}

// String returns a human-readable description of this Scoped View ID
func (id ScopedViewId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("View Name: %q", id.ViewName),
	}
	return fmt.Sprintf("Scoped View (%s)", strings.Join(components, "\n"))
}
