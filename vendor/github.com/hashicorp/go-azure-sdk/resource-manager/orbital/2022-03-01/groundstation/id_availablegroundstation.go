package groundstation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AvailableGroundStationId{}

// AvailableGroundStationId is a struct representing the Resource ID for a Available Ground Station
type AvailableGroundStationId struct {
	SubscriptionId             string
	AvailableGroundStationName string
}

// NewAvailableGroundStationID returns a new AvailableGroundStationId struct
func NewAvailableGroundStationID(subscriptionId string, availableGroundStationName string) AvailableGroundStationId {
	return AvailableGroundStationId{
		SubscriptionId:             subscriptionId,
		AvailableGroundStationName: availableGroundStationName,
	}
}

// ParseAvailableGroundStationID parses 'input' into a AvailableGroundStationId
func ParseAvailableGroundStationID(input string) (*AvailableGroundStationId, error) {
	parser := resourceids.NewParserFromResourceIdType(AvailableGroundStationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AvailableGroundStationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.AvailableGroundStationName, ok = parsed.Parsed["availableGroundStationName"]; !ok {
		return nil, fmt.Errorf("the segment 'availableGroundStationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAvailableGroundStationIDInsensitively parses 'input' case-insensitively into a AvailableGroundStationId
// note: this method should only be used for API response data and not user input
func ParseAvailableGroundStationIDInsensitively(input string) (*AvailableGroundStationId, error) {
	parser := resourceids.NewParserFromResourceIdType(AvailableGroundStationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AvailableGroundStationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.AvailableGroundStationName, ok = parsed.Parsed["availableGroundStationName"]; !ok {
		return nil, fmt.Errorf("the segment 'availableGroundStationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAvailableGroundStationID checks that 'input' can be parsed as a Available Ground Station ID
func ValidateAvailableGroundStationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAvailableGroundStationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Available Ground Station ID
func (id AvailableGroundStationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Orbital/availableGroundStations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AvailableGroundStationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Available Ground Station ID
func (id AvailableGroundStationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOrbital", "Microsoft.Orbital", "Microsoft.Orbital"),
		resourceids.StaticSegment("staticAvailableGroundStations", "availableGroundStations", "availableGroundStations"),
		resourceids.UserSpecifiedSegment("availableGroundStationName", "availableGroundStationValue"),
	}
}

// String returns a human-readable description of this Available Ground Station ID
func (id AvailableGroundStationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Available Ground Station Name: %q", id.AvailableGroundStationName),
	}
	return fmt.Sprintf("Available Ground Station (%s)", strings.Join(components, "\n"))
}
