package raicontentfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RaiContentFilterId{})
}

var _ resourceids.ResourceId = &RaiContentFilterId{}

// RaiContentFilterId is a struct representing the Resource ID for a Rai Content Filter
type RaiContentFilterId struct {
	SubscriptionId       string
	LocationName         string
	RaiContentFilterName string
}

// NewRaiContentFilterID returns a new RaiContentFilterId struct
func NewRaiContentFilterID(subscriptionId string, locationName string, raiContentFilterName string) RaiContentFilterId {
	return RaiContentFilterId{
		SubscriptionId:       subscriptionId,
		LocationName:         locationName,
		RaiContentFilterName: raiContentFilterName,
	}
}

// ParseRaiContentFilterID parses 'input' into a RaiContentFilterId
func ParseRaiContentFilterID(input string) (*RaiContentFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiContentFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiContentFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRaiContentFilterIDInsensitively parses 'input' case-insensitively into a RaiContentFilterId
// note: this method should only be used for API response data and not user input
func ParseRaiContentFilterIDInsensitively(input string) (*RaiContentFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiContentFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiContentFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RaiContentFilterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.RaiContentFilterName, ok = input.Parsed["raiContentFilterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "raiContentFilterName", input)
	}

	return nil
}

// ValidateRaiContentFilterID checks that 'input' can be parsed as a Rai Content Filter ID
func ValidateRaiContentFilterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRaiContentFilterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rai Content Filter ID
func (id RaiContentFilterId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CognitiveServices/locations/%s/raiContentFilters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.RaiContentFilterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rai Content Filter ID
func (id RaiContentFilterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticRaiContentFilters", "raiContentFilters", "raiContentFilters"),
		resourceids.UserSpecifiedSegment("raiContentFilterName", "raiContentFilterName"),
	}
}

// String returns a human-readable description of this Rai Content Filter ID
func (id RaiContentFilterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Rai Content Filter Name: %q", id.RaiContentFilterName),
	}
	return fmt.Sprintf("Rai Content Filter (%s)", strings.Join(components, "\n"))
}
