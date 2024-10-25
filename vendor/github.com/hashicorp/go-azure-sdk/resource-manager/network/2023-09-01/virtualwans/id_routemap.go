package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteMapId{})
}

var _ resourceids.ResourceId = &RouteMapId{}

// RouteMapId is a struct representing the Resource ID for a Route Map
type RouteMapId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualHubName    string
	RouteMapName      string
}

// NewRouteMapID returns a new RouteMapId struct
func NewRouteMapID(subscriptionId string, resourceGroupName string, virtualHubName string, routeMapName string) RouteMapId {
	return RouteMapId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualHubName:    virtualHubName,
		RouteMapName:      routeMapName,
	}
}

// ParseRouteMapID parses 'input' into a RouteMapId
func ParseRouteMapID(input string) (*RouteMapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteMapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteMapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteMapIDInsensitively parses 'input' case-insensitively into a RouteMapId
// note: this method should only be used for API response data and not user input
func ParseRouteMapIDInsensitively(input string) (*RouteMapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteMapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteMapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteMapId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualHubName, ok = input.Parsed["virtualHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", input)
	}

	if id.RouteMapName, ok = input.Parsed["routeMapName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeMapName", input)
	}

	return nil
}

// ValidateRouteMapID checks that 'input' can be parsed as a Route Map ID
func ValidateRouteMapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteMapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Map ID
func (id RouteMapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/routeMaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.RouteMapName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Map ID
func (id RouteMapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticRouteMaps", "routeMaps", "routeMaps"),
		resourceids.UserSpecifiedSegment("routeMapName", "routeMapName"),
	}
}

// String returns a human-readable description of this Route Map ID
func (id RouteMapId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Route Map Name: %q", id.RouteMapName),
	}
	return fmt.Sprintf("Route Map (%s)", strings.Join(components, "\n"))
}
