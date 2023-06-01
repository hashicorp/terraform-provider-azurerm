package exports

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedExportId{}

// ScopedExportId is a struct representing the Resource ID for a Scoped Export
type ScopedExportId struct {
	Scope      string
	ExportName string
}

// NewScopedExportID returns a new ScopedExportId struct
func NewScopedExportID(scope string, exportName string) ScopedExportId {
	return ScopedExportId{
		Scope:      scope,
		ExportName: exportName,
	}
}

// ParseScopedExportID parses 'input' into a ScopedExportId
func ParseScopedExportID(input string) (*ScopedExportId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedExportId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedExportId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ExportName, ok = parsed.Parsed["exportName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "exportName", *parsed)
	}

	return &id, nil
}

// ParseScopedExportIDInsensitively parses 'input' case-insensitively into a ScopedExportId
// note: this method should only be used for API response data and not user input
func ParseScopedExportIDInsensitively(input string) (*ScopedExportId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedExportId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedExportId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.ExportName, ok = parsed.Parsed["exportName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "exportName", *parsed)
	}

	return &id, nil
}

// ValidateScopedExportID checks that 'input' can be parsed as a Scoped Export ID
func ValidateScopedExportID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedExportID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Export ID
func (id ScopedExportId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.ExportName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Export ID
func (id ScopedExportId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticExports", "exports", "exports"),
		resourceids.UserSpecifiedSegment("exportName", "exportValue"),
	}
}

// String returns a human-readable description of this Scoped Export ID
func (id ScopedExportId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Export Name: %q", id.ExportName),
	}
	return fmt.Sprintf("Scoped Export (%s)", strings.Join(components, "\n"))
}
