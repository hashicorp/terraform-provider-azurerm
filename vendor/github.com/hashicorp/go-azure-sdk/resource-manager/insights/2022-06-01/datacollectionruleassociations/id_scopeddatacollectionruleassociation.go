package datacollectionruleassociations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedDataCollectionRuleAssociationId{}

// ScopedDataCollectionRuleAssociationId is a struct representing the Resource ID for a Scoped Data Collection Rule Association
type ScopedDataCollectionRuleAssociationId struct {
	ResourceUri                       string
	DataCollectionRuleAssociationName string
}

// NewScopedDataCollectionRuleAssociationID returns a new ScopedDataCollectionRuleAssociationId struct
func NewScopedDataCollectionRuleAssociationID(resourceUri string, dataCollectionRuleAssociationName string) ScopedDataCollectionRuleAssociationId {
	return ScopedDataCollectionRuleAssociationId{
		ResourceUri:                       resourceUri,
		DataCollectionRuleAssociationName: dataCollectionRuleAssociationName,
	}
}

// ParseScopedDataCollectionRuleAssociationID parses 'input' into a ScopedDataCollectionRuleAssociationId
func ParseScopedDataCollectionRuleAssociationID(input string) (*ScopedDataCollectionRuleAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedDataCollectionRuleAssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedDataCollectionRuleAssociationId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.DataCollectionRuleAssociationName, ok = parsed.Parsed["dataCollectionRuleAssociationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataCollectionRuleAssociationName", *parsed)
	}

	return &id, nil
}

// ParseScopedDataCollectionRuleAssociationIDInsensitively parses 'input' case-insensitively into a ScopedDataCollectionRuleAssociationId
// note: this method should only be used for API response data and not user input
func ParseScopedDataCollectionRuleAssociationIDInsensitively(input string) (*ScopedDataCollectionRuleAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedDataCollectionRuleAssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedDataCollectionRuleAssociationId{}

	if id.ResourceUri, ok = parsed.Parsed["resourceUri"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", *parsed)
	}

	if id.DataCollectionRuleAssociationName, ok = parsed.Parsed["dataCollectionRuleAssociationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataCollectionRuleAssociationName", *parsed)
	}

	return &id, nil
}

// ValidateScopedDataCollectionRuleAssociationID checks that 'input' can be parsed as a Scoped Data Collection Rule Association ID
func ValidateScopedDataCollectionRuleAssociationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedDataCollectionRuleAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Data Collection Rule Association ID
func (id ScopedDataCollectionRuleAssociationId) ID() string {
	fmtString := "/%s/providers/Microsoft.Insights/dataCollectionRuleAssociations/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.DataCollectionRuleAssociationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Data Collection Rule Association ID
func (id ScopedDataCollectionRuleAssociationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDataCollectionRuleAssociations", "dataCollectionRuleAssociations", "dataCollectionRuleAssociations"),
		resourceids.UserSpecifiedSegment("dataCollectionRuleAssociationName", "dataCollectionRuleAssociationValue"),
	}
}

// String returns a human-readable description of this Scoped Data Collection Rule Association ID
func (id ScopedDataCollectionRuleAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Data Collection Rule Association Name: %q", id.DataCollectionRuleAssociationName),
	}
	return fmt.Sprintf("Scoped Data Collection Rule Association (%s)", strings.Join(components, "\n"))
}
