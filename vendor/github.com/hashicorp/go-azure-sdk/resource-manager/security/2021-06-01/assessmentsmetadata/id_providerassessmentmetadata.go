package assessmentsmetadata

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderAssessmentMetadataId{}

// ProviderAssessmentMetadataId is a struct representing the Resource ID for a Provider Assessment Metadata
type ProviderAssessmentMetadataId struct {
	SubscriptionId         string
	AssessmentMetadataName string
}

// NewProviderAssessmentMetadataID returns a new ProviderAssessmentMetadataId struct
func NewProviderAssessmentMetadataID(subscriptionId string, assessmentMetadataName string) ProviderAssessmentMetadataId {
	return ProviderAssessmentMetadataId{
		SubscriptionId:         subscriptionId,
		AssessmentMetadataName: assessmentMetadataName,
	}
}

// ParseProviderAssessmentMetadataID parses 'input' into a ProviderAssessmentMetadataId
func ParseProviderAssessmentMetadataID(input string) (*ProviderAssessmentMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderAssessmentMetadataId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderAssessmentMetadataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.AssessmentMetadataName, ok = parsed.Parsed["assessmentMetadataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assessmentMetadataName", *parsed)
	}

	return &id, nil
}

// ParseProviderAssessmentMetadataIDInsensitively parses 'input' case-insensitively into a ProviderAssessmentMetadataId
// note: this method should only be used for API response data and not user input
func ParseProviderAssessmentMetadataIDInsensitively(input string) (*ProviderAssessmentMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderAssessmentMetadataId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderAssessmentMetadataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.AssessmentMetadataName, ok = parsed.Parsed["assessmentMetadataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assessmentMetadataName", *parsed)
	}

	return &id, nil
}

// ValidateProviderAssessmentMetadataID checks that 'input' can be parsed as a Provider Assessment Metadata ID
func ValidateProviderAssessmentMetadataID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderAssessmentMetadataID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Assessment Metadata ID
func (id ProviderAssessmentMetadataId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/assessmentMetadata/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AssessmentMetadataName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Assessment Metadata ID
func (id ProviderAssessmentMetadataId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurity", "Microsoft.Security", "Microsoft.Security"),
		resourceids.StaticSegment("staticAssessmentMetadata", "assessmentMetadata", "assessmentMetadata"),
		resourceids.UserSpecifiedSegment("assessmentMetadataName", "assessmentMetadataValue"),
	}
}

// String returns a human-readable description of this Provider Assessment Metadata ID
func (id ProviderAssessmentMetadataId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Assessment Metadata Name: %q", id.AssessmentMetadataName),
	}
	return fmt.Sprintf("Provider Assessment Metadata (%s)", strings.Join(components, "\n"))
}
