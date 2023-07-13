package diagnosticsettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedDiagnosticSettingId{}

// ScopedDiagnosticSettingId is a struct representing the Resource ID for a Scoped Diagnostic Setting
type ScopedDiagnosticSettingId struct {
	ResourceUri           string
	DiagnosticSettingName string
}

// NewScopedDiagnosticSettingID returns a new ScopedDiagnosticSettingId struct
func NewScopedDiagnosticSettingID(resourceUri string, diagnosticSettingName string) ScopedDiagnosticSettingId {
	return ScopedDiagnosticSettingId{
		ResourceUri:           resourceUri,
		DiagnosticSettingName: diagnosticSettingName,
	}
}

// ParseScopedDiagnosticSettingID parses 'input' into a ScopedDiagnosticSettingId
func ParseScopedDiagnosticSettingID(input string) (*ScopedDiagnosticSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedDiagnosticSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedDiagnosticSettingId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.DiagnosticSettingName, ok = parsed.Parsed["diagnosticSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diagnosticSettingName", *parsed)
	}

	return &id, nil
}

// ParseScopedDiagnosticSettingIDInsensitively parses 'input' case-insensitively into a ScopedDiagnosticSettingId
// note: this method should only be used for API response data and not user input
func ParseScopedDiagnosticSettingIDInsensitively(input string) (*ScopedDiagnosticSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedDiagnosticSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedDiagnosticSettingId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.DiagnosticSettingName, ok = parsed.Parsed["diagnosticSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diagnosticSettingName", *parsed)
	}

	return &id, nil
}

// ValidateScopedDiagnosticSettingID checks that 'input' can be parsed as a Scoped Diagnostic Setting ID
func ValidateScopedDiagnosticSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedDiagnosticSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Diagnostic Setting ID
func (id ScopedDiagnosticSettingId) ID() string {
	fmtString := "/%s/providers/Microsoft.Insights/diagnosticSettings/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.DiagnosticSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Diagnostic Setting ID
func (id ScopedDiagnosticSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDiagnosticSettings", "diagnosticSettings", "diagnosticSettings"),
		resourceids.UserSpecifiedSegment("diagnosticSettingName", "diagnosticSettingValue"),
	}
}

// String returns a human-readable description of this Scoped Diagnostic Setting ID
func (id ScopedDiagnosticSettingId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Diagnostic Setting Name: %q", id.DiagnosticSettingName),
	}
	return fmt.Sprintf("Scoped Diagnostic Setting (%s)", strings.Join(components, "\n"))
}
