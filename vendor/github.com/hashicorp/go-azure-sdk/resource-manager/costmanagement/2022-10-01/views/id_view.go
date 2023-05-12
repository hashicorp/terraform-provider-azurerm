package views

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ViewId{}

// ViewId is a struct representing the Resource ID for a View
type ViewId struct {
	ViewName string
}

// NewViewID returns a new ViewId struct
func NewViewID(viewName string) ViewId {
	return ViewId{
		ViewName: viewName,
	}
}

// ParseViewID parses 'input' into a ViewId
func ParseViewID(input string) (*ViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(ViewId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ViewId{}

	if id.ViewName, ok = parsed.Parsed["viewName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "viewName", *parsed)
	}

	return &id, nil
}

// ParseViewIDInsensitively parses 'input' case-insensitively into a ViewId
// note: this method should only be used for API response data and not user input
func ParseViewIDInsensitively(input string) (*ViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(ViewId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ViewId{}

	if id.ViewName, ok = parsed.Parsed["viewName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "viewName", *parsed)
	}

	return &id, nil
}

// ValidateViewID checks that 'input' can be parsed as a View ID
func ValidateViewID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseViewID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted View ID
func (id ViewId) ID() string {
	fmtString := "/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, id.ViewName)
}

// Segments returns a slice of Resource ID Segments which comprise this View ID
func (id ViewId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticViews", "views", "views"),
		resourceids.UserSpecifiedSegment("viewName", "viewValue"),
	}
}

// String returns a human-readable description of this View ID
func (id ViewId) String() string {
	components := []string{
		fmt.Sprintf("View Name: %q", id.ViewName),
	}
	return fmt.Sprintf("View (%s)", strings.Join(components, "\n"))
}
