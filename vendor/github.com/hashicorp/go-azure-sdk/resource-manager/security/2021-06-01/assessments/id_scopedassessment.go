package assessments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedAssessmentId{}

// ScopedAssessmentId is a struct representing the Resource ID for a Scoped Assessment
type ScopedAssessmentId struct {
	ResourceId     string
	AssessmentName string
}

// NewScopedAssessmentID returns a new ScopedAssessmentId struct
func NewScopedAssessmentID(resourceId string, assessmentName string) ScopedAssessmentId {
	return ScopedAssessmentId{
		ResourceId:     resourceId,
		AssessmentName: assessmentName,
	}
}

// ParseScopedAssessmentID parses 'input' into a ScopedAssessmentId
func ParseScopedAssessmentID(input string) (*ScopedAssessmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedAssessmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedAssessmentId{}

	if id.ResourceId, ok = parsed.Parsed["resourceId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceId", *parsed)
	}

	if id.AssessmentName, ok = parsed.Parsed["assessmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assessmentName", *parsed)
	}

	return &id, nil
}

// ParseScopedAssessmentIDInsensitively parses 'input' case-insensitively into a ScopedAssessmentId
// note: this method should only be used for API response data and not user input
func ParseScopedAssessmentIDInsensitively(input string) (*ScopedAssessmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedAssessmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedAssessmentId{}

	if id.ResourceId, ok = parsed.Parsed["resourceId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceId", *parsed)
	}

	if id.AssessmentName, ok = parsed.Parsed["assessmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assessmentName", *parsed)
	}

	return &id, nil
}

// ValidateScopedAssessmentID checks that 'input' can be parsed as a Scoped Assessment ID
func ValidateScopedAssessmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedAssessmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Assessment ID
func (id ScopedAssessmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Security/assessments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceId, "/"), id.AssessmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Assessment ID
func (id ScopedAssessmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurity", "Microsoft.Security", "Microsoft.Security"),
		resourceids.StaticSegment("staticAssessments", "assessments", "assessments"),
		resourceids.UserSpecifiedSegment("assessmentName", "assessmentValue"),
	}
}

// String returns a human-readable description of this Scoped Assessment ID
func (id ScopedAssessmentId) String() string {
	components := []string{
		fmt.Sprintf("Resource: %q", id.ResourceId),
		fmt.Sprintf("Assessment Name: %q", id.AssessmentName),
	}
	return fmt.Sprintf("Scoped Assessment (%s)", strings.Join(components, "\n"))
}
