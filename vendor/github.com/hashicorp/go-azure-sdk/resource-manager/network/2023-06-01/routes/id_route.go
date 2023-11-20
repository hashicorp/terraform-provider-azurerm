package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RouteId{}

// RouteId is a struct representing the Resource ID for a Route
type RouteId struct {
	SubscriptionId    string
	ResourceGroupName string
	RouteTableName    string
	RouteName         string
}

// NewRouteID returns a new RouteId struct
func NewRouteID(subscriptionId string, resourceGroupName string, routeTableName string, routeName string) RouteId {
	return RouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RouteTableName:    routeTableName,
		RouteName:         routeName,
	}
}

// ParseRouteID parses 'input' into a RouteId
func ParseRouteID(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RouteTableName, ok = parsed.Parsed["routeTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", *parsed)
	}

	if id.RouteName, ok = parsed.Parsed["routeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeName", *parsed)
	}

	return &id, nil
}

// ParseRouteIDInsensitively parses 'input' case-insensitively into a RouteId
// note: this method should only be used for API response data and not user input
func ParseRouteIDInsensitively(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RouteTableName, ok = parsed.Parsed["routeTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", *parsed)
	}

	if id.RouteName, ok = parsed.Parsed["routeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeName", *parsed)
	}

	return &id, nil
}

// ValidateRouteID checks that 'input' can be parsed as a Route ID
func ValidateRouteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route ID
func (id RouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RouteTableName, id.RouteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route ID
func (id RouteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableValue"),
		resourceids.StaticSegment("staticRoutes", "routes", "routes"),
		resourceids.UserSpecifiedSegment("routeName", "routeValue"),
	}
}

// String returns a human-readable description of this Route ID
func (id RouteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Route Table Name: %q", id.RouteTableName),
		fmt.Sprintf("Route Name: %q", id.RouteName),
	}
	return fmt.Sprintf("Route (%s)", strings.Join(components, "\n"))
}
