package routefilterrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteFilterId{})
}

var _ resourceids.ResourceId = &RouteFilterId{}

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
	parser := resourceids.NewParserFromResourceIdType(&RouteFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteFilterIDInsensitively parses 'input' case-insensitively into a RouteFilterId
// note: this method should only be used for API response data and not user input
func ParseRouteFilterIDInsensitively(input string) (*RouteFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteFilterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RouteFilterName, ok = input.Parsed["routeFilterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeFilterName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("routeFilterName", "routeFilterName"),
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
