package diagnosticsettingscategories

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedDiagnosticSettingsCategoryId{})
}

var _ resourceids.ResourceId = &ScopedDiagnosticSettingsCategoryId{}

// ScopedDiagnosticSettingsCategoryId is a struct representing the Resource ID for a Scoped Diagnostic Settings Category
type ScopedDiagnosticSettingsCategoryId struct {
	ResourceUri                    string
	DiagnosticSettingsCategoryName string
}

// NewScopedDiagnosticSettingsCategoryID returns a new ScopedDiagnosticSettingsCategoryId struct
func NewScopedDiagnosticSettingsCategoryID(resourceUri string, diagnosticSettingsCategoryName string) ScopedDiagnosticSettingsCategoryId {
	return ScopedDiagnosticSettingsCategoryId{
		ResourceUri:                    resourceUri,
		DiagnosticSettingsCategoryName: diagnosticSettingsCategoryName,
	}
}

// ParseScopedDiagnosticSettingsCategoryID parses 'input' into a ScopedDiagnosticSettingsCategoryId
func ParseScopedDiagnosticSettingsCategoryID(input string) (*ScopedDiagnosticSettingsCategoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedDiagnosticSettingsCategoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedDiagnosticSettingsCategoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedDiagnosticSettingsCategoryIDInsensitively parses 'input' case-insensitively into a ScopedDiagnosticSettingsCategoryId
// note: this method should only be used for API response data and not user input
func ParseScopedDiagnosticSettingsCategoryIDInsensitively(input string) (*ScopedDiagnosticSettingsCategoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedDiagnosticSettingsCategoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedDiagnosticSettingsCategoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedDiagnosticSettingsCategoryId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceUri, ok = input.Parsed["resourceUri"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", input)
	}

	if id.DiagnosticSettingsCategoryName, ok = input.Parsed["diagnosticSettingsCategoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diagnosticSettingsCategoryName", input)
	}

	return nil
}

// ValidateScopedDiagnosticSettingsCategoryID checks that 'input' can be parsed as a Scoped Diagnostic Settings Category ID
func ValidateScopedDiagnosticSettingsCategoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedDiagnosticSettingsCategoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Diagnostic Settings Category ID
func (id ScopedDiagnosticSettingsCategoryId) ID() string {
	fmtString := "/%s/providers/Microsoft.Insights/diagnosticSettingsCategories/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.DiagnosticSettingsCategoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Diagnostic Settings Category ID
func (id ScopedDiagnosticSettingsCategoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDiagnosticSettingsCategories", "diagnosticSettingsCategories", "diagnosticSettingsCategories"),
		resourceids.UserSpecifiedSegment("diagnosticSettingsCategoryName", "diagnosticSettingsCategoryName"),
	}
}

// String returns a human-readable description of this Scoped Diagnostic Settings Category ID
func (id ScopedDiagnosticSettingsCategoryId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Diagnostic Settings Category Name: %q", id.DiagnosticSettingsCategoryName),
	}
	return fmt.Sprintf("Scoped Diagnostic Settings Category (%s)", strings.Join(components, "\n"))
}
