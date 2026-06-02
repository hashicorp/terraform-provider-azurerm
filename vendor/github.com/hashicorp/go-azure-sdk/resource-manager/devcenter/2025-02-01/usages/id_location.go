package usages

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocationId{})
}

var _ resourceids.ResourceId = &LocationId{}

// LocationId is a struct representing the Resource ID for a Location
type LocationId struct {
	SubscriptionId string
	LocationName   string
}

// NewLocationID returns a new LocationId struct
func NewLocationID(subscriptionId string, locationName string) LocationId {
	return LocationId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
	}
}

// ParseLocationID parses 'input' into a LocationId
func ParseLocationID(input string) (*LocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocationIDInsensitively parses 'input' case-insensitively into a LocationId
// note: this method should only be used for API response data and not user input
func ParseLocationIDInsensitively(input string) (*LocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	return nil
}

// ValidateLocationID checks that 'input' can be parsed as a Location ID
func ValidateLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location ID
func (id LocationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DevCenter/locations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Location ID
func (id LocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
	}
}

// String returns a human-readable description of this Location ID
func (id LocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
	}
	return fmt.Sprintf("Location (%s)", strings.Join(components, "\n"))
}
