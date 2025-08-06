package suppressions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedSuppressionId{})
}

var _ resourceids.ResourceId = &ScopedSuppressionId{}

// ScopedSuppressionId is a struct representing the Resource ID for a Scoped Suppression
type ScopedSuppressionId struct {
	ResourceUri      string
	RecommendationId string
	SuppressionName  string
}

// NewScopedSuppressionID returns a new ScopedSuppressionId struct
func NewScopedSuppressionID(resourceUri string, recommendationId string, suppressionName string) ScopedSuppressionId {
	return ScopedSuppressionId{
		ResourceUri:      resourceUri,
		RecommendationId: recommendationId,
		SuppressionName:  suppressionName,
	}
}

// ParseScopedSuppressionID parses 'input' into a ScopedSuppressionId
func ParseScopedSuppressionID(input string) (*ScopedSuppressionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedSuppressionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedSuppressionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedSuppressionIDInsensitively parses 'input' case-insensitively into a ScopedSuppressionId
// note: this method should only be used for API response data and not user input
func ParseScopedSuppressionIDInsensitively(input string) (*ScopedSuppressionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedSuppressionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedSuppressionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedSuppressionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceUri, ok = input.Parsed["resourceUri"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", input)
	}

	if id.RecommendationId, ok = input.Parsed["recommendationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "recommendationId", input)
	}

	if id.SuppressionName, ok = input.Parsed["suppressionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "suppressionName", input)
	}

	return nil
}

// ValidateScopedSuppressionID checks that 'input' can be parsed as a Scoped Suppression ID
func ValidateScopedSuppressionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedSuppressionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Suppression ID
func (id ScopedSuppressionId) ID() string {
	fmtString := "/%s/providers/Microsoft.Advisor/recommendations/%s/suppressions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.RecommendationId, id.SuppressionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Suppression ID
func (id ScopedSuppressionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAdvisor", "Microsoft.Advisor", "Microsoft.Advisor"),
		resourceids.StaticSegment("staticRecommendations", "recommendations", "recommendations"),
		resourceids.UserSpecifiedSegment("recommendationId", "recommendationId"),
		resourceids.StaticSegment("staticSuppressions", "suppressions", "suppressions"),
		resourceids.UserSpecifiedSegment("suppressionName", "suppressionName"),
	}
}

// String returns a human-readable description of this Scoped Suppression ID
func (id ScopedSuppressionId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Recommendation: %q", id.RecommendationId),
		fmt.Sprintf("Suppression Name: %q", id.SuppressionName),
	}
	return fmt.Sprintf("Scoped Suppression (%s)", strings.Join(components, "\n"))
}
