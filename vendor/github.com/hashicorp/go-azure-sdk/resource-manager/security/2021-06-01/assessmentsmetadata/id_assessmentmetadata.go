package assessmentsmetadata

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AssessmentMetadataId{})
}

var _ resourceids.ResourceId = &AssessmentMetadataId{}

// AssessmentMetadataId is a struct representing the Resource ID for a Assessment Metadata
type AssessmentMetadataId struct {
	AssessmentMetadataName string
}

// NewAssessmentMetadataID returns a new AssessmentMetadataId struct
func NewAssessmentMetadataID(assessmentMetadataName string) AssessmentMetadataId {
	return AssessmentMetadataId{
		AssessmentMetadataName: assessmentMetadataName,
	}
}

// ParseAssessmentMetadataID parses 'input' into a AssessmentMetadataId
func ParseAssessmentMetadataID(input string) (*AssessmentMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssessmentMetadataId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssessmentMetadataId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAssessmentMetadataIDInsensitively parses 'input' case-insensitively into a AssessmentMetadataId
// note: this method should only be used for API response data and not user input
func ParseAssessmentMetadataIDInsensitively(input string) (*AssessmentMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssessmentMetadataId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssessmentMetadataId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AssessmentMetadataId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.AssessmentMetadataName, ok = input.Parsed["assessmentMetadataName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "assessmentMetadataName", input)
	}

	return nil
}

// ValidateAssessmentMetadataID checks that 'input' can be parsed as a Assessment Metadata ID
func ValidateAssessmentMetadataID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssessmentMetadataID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Assessment Metadata ID
func (id AssessmentMetadataId) ID() string {
	fmtString := "/providers/Microsoft.Security/assessmentMetadata/%s"
	return fmt.Sprintf(fmtString, id.AssessmentMetadataName)
}

// Segments returns a slice of Resource ID Segments which comprise this Assessment Metadata ID
func (id AssessmentMetadataId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurity", "Microsoft.Security", "Microsoft.Security"),
		resourceids.StaticSegment("staticAssessmentMetadata", "assessmentMetadata", "assessmentMetadata"),
		resourceids.UserSpecifiedSegment("assessmentMetadataName", "assessmentMetadataName"),
	}
}

// String returns a human-readable description of this Assessment Metadata ID
func (id AssessmentMetadataId) String() string {
	components := []string{
		fmt.Sprintf("Assessment Metadata Name: %q", id.AssessmentMetadataName),
	}
	return fmt.Sprintf("Assessment Metadata (%s)", strings.Join(components, "\n"))
}
