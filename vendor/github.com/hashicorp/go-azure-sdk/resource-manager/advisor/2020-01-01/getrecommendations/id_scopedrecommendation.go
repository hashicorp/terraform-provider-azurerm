package getrecommendations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedRecommendationId{}

// ScopedRecommendationId is a struct representing the Resource ID for a Scoped Recommendation
type ScopedRecommendationId struct {
	ResourceUri      string
	RecommendationId string
}

// NewScopedRecommendationID returns a new ScopedRecommendationId struct
func NewScopedRecommendationID(resourceUri string, recommendationId string) ScopedRecommendationId {
	return ScopedRecommendationId{
		ResourceUri:      resourceUri,
		RecommendationId: recommendationId,
	}
}

// ParseScopedRecommendationID parses 'input' into a ScopedRecommendationId
func ParseScopedRecommendationID(input string) (*ScopedRecommendationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRecommendationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRecommendationId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.RecommendationId, ok = parsed.Parsed["recommendationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "recommendationId", *parsed)
	}

	return &id, nil
}

// ParseScopedRecommendationIDInsensitively parses 'input' case-insensitively into a ScopedRecommendationId
// note: this method should only be used for API response data and not user input
func ParseScopedRecommendationIDInsensitively(input string) (*ScopedRecommendationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRecommendationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRecommendationId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.RecommendationId, ok = parsed.Parsed["recommendationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "recommendationId", *parsed)
	}

	return &id, nil
}

// ValidateScopedRecommendationID checks that 'input' can be parsed as a Scoped Recommendation ID
func ValidateScopedRecommendationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRecommendationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Recommendation ID
func (id ScopedRecommendationId) ID() string {
	fmtString := "/%s/providers/Microsoft.Advisor/recommendations/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.RecommendationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Recommendation ID
func (id ScopedRecommendationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAdvisor", "Microsoft.Advisor", "Microsoft.Advisor"),
		resourceids.StaticSegment("staticRecommendations", "recommendations", "recommendations"),
		resourceids.UserSpecifiedSegment("recommendationId", "recommendationIdValue"),
	}
}

// String returns a human-readable description of this Scoped Recommendation ID
func (id ScopedRecommendationId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Recommendation: %q", id.RecommendationId),
	}
	return fmt.Sprintf("Scoped Recommendation (%s)", strings.Join(components, "\n"))
}
