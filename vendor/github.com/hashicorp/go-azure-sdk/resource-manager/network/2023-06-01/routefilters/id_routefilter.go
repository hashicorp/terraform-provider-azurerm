package routefilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RouteFilterId{}

// RouteFilterId is a struct representing the Resource ID for a Route Filter
type RouteFilterId struct {
	SubscriptionId    string
	ResourceGroupName string
	RouteFilterName   string
}

// NewRouteFilterID returns a new RouteFilterId struct
func NewRouteFilterID(subscriptionId string, resourceGroupName string, routeFilterName string) RouteFilterId {
	return RouteFilterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RouteFilterName:   routeFilterName,
	}
}

// ParseRouteFilterID parses 'input' into a RouteFilterId
func ParseRouteFilterID(input string) (*RouteFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RouteFilterName, ok = parsed.Parsed["routeFilterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeFilterName", *parsed)
	}

	return &id, nil
}

// ParseRouteFilterIDInsensitively parses 'input' case-insensitively into a RouteFilterId
// note: this method should only be used for API response data and not user input
func ParseRouteFilterIDInsensitively(input string) (*RouteFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RouteFilterName, ok = parsed.Parsed["routeFilterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeFilterName", *parsed)
	}

	return &id, nil
}

// ValidateRouteFilterID checks that 'input' can be parsed as a Route Filter ID
func ValidateRouteFilterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteFilterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Filter ID
func (id RouteFilterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeFilters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RouteFilterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Filter ID
func (id RouteFilterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticRouteFilters", "routeFilters", "routeFilters"),
		resourceids.UserSpecifiedSegment("routeFilterName", "routeFilterValue"),
	}
}

// String returns a human-readable description of this Route Filter ID
func (id RouteFilterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Route Filter Name: %q", id.RouteFilterName),
	}
	return fmt.Sprintf("Route Filter (%s)", strings.Join(components, "\n"))
}
