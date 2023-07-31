package expressrouteportslocations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExpressRoutePortsLocationId{}

// ExpressRoutePortsLocationId is a struct representing the Resource ID for a Express Route Ports Location
type ExpressRoutePortsLocationId struct {
	SubscriptionId                string
	ExpressRoutePortsLocationName string
}

// NewExpressRoutePortsLocationID returns a new ExpressRoutePortsLocationId struct
func NewExpressRoutePortsLocationID(subscriptionId string, expressRoutePortsLocationName string) ExpressRoutePortsLocationId {
	return ExpressRoutePortsLocationId{
		SubscriptionId:                subscriptionId,
		ExpressRoutePortsLocationName: expressRoutePortsLocationName,
	}
}

// ParseExpressRoutePortsLocationID parses 'input' into a ExpressRoutePortsLocationId
func ParseExpressRoutePortsLocationID(input string) (*ExpressRoutePortsLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRoutePortsLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRoutePortsLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ExpressRoutePortsLocationName, ok = parsed.Parsed["expressRoutePortsLocationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortsLocationName", *parsed)
	}

	return &id, nil
}

// ParseExpressRoutePortsLocationIDInsensitively parses 'input' case-insensitively into a ExpressRoutePortsLocationId
// note: this method should only be used for API response data and not user input
func ParseExpressRoutePortsLocationIDInsensitively(input string) (*ExpressRoutePortsLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRoutePortsLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRoutePortsLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ExpressRoutePortsLocationName, ok = parsed.Parsed["expressRoutePortsLocationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortsLocationName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRoutePortsLocationID checks that 'input' can be parsed as a Express Route Ports Location ID
func ValidateExpressRoutePortsLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRoutePortsLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Ports Location ID
func (id ExpressRoutePortsLocationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/expressRoutePortsLocations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ExpressRoutePortsLocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Ports Location ID
func (id ExpressRoutePortsLocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRoutePortsLocations", "expressRoutePortsLocations", "expressRoutePortsLocations"),
		resourceids.UserSpecifiedSegment("expressRoutePortsLocationName", "expressRoutePortsLocationValue"),
	}
}

// String returns a human-readable description of this Express Route Ports Location ID
func (id ExpressRoutePortsLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Express Route Ports Location Name: %q", id.ExpressRoutePortsLocationName),
	}
	return fmt.Sprintf("Express Route Ports Location (%s)", strings.Join(components, "\n"))
}
