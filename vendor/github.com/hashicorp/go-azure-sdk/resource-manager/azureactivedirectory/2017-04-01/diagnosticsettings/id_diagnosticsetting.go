package diagnosticsettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DiagnosticSettingId{}

// DiagnosticSettingId is a struct representing the Resource ID for a Diagnostic Setting
type DiagnosticSettingId struct {
	DiagnosticSettingName string
}

// NewDiagnosticSettingID returns a new DiagnosticSettingId struct
func NewDiagnosticSettingID(diagnosticSettingName string) DiagnosticSettingId {
	return DiagnosticSettingId{
		DiagnosticSettingName: diagnosticSettingName,
	}
}

// ParseDiagnosticSettingID parses 'input' into a DiagnosticSettingId
func ParseDiagnosticSettingID(input string) (*DiagnosticSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiagnosticSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiagnosticSettingId{}

	if id.DiagnosticSettingName, ok = parsed.Parsed["diagnosticSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diagnosticSettingName", *parsed)
	}

	return &id, nil
}

// ParseDiagnosticSettingIDInsensitively parses 'input' case-insensitively into a DiagnosticSettingId
// note: this method should only be used for API response data and not user input
func ParseDiagnosticSettingIDInsensitively(input string) (*DiagnosticSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiagnosticSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiagnosticSettingId{}

	if id.DiagnosticSettingName, ok = parsed.Parsed["diagnosticSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diagnosticSettingName", *parsed)
	}

	return &id, nil
}

// ValidateDiagnosticSettingID checks that 'input' can be parsed as a Diagnostic Setting ID
func ValidateDiagnosticSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDiagnosticSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Diagnostic Setting ID
func (id DiagnosticSettingId) ID() string {
	fmtString := "/providers/Microsoft.AADIAM/diagnosticSettings/%s"
	return fmt.Sprintf(fmtString, id.DiagnosticSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Diagnostic Setting ID
func (id DiagnosticSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAADIAM", "Microsoft.AADIAM", "Microsoft.AADIAM"),
		resourceids.StaticSegment("staticDiagnosticSettings", "diagnosticSettings", "diagnosticSettings"),
		resourceids.UserSpecifiedSegment("diagnosticSettingName", "diagnosticSettingValue"),
	}
}

// String returns a human-readable description of this Diagnostic Setting ID
func (id DiagnosticSettingId) String() string {
	components := []string{
		fmt.Sprintf("Diagnostic Setting Name: %q", id.DiagnosticSettingName),
	}
	return fmt.Sprintf("Diagnostic Setting (%s)", strings.Join(components, "\n"))
}
