package expressroutecircuitroutestablesummary

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RouteTablesSummaryId{}

// RouteTablesSummaryId is a struct representing the Resource ID for a Route Tables Summary
type RouteTablesSummaryId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
	PeeringName             string
	RouteTablesSummaryName  string
}

// NewRouteTablesSummaryID returns a new RouteTablesSummaryId struct
func NewRouteTablesSummaryID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string, peeringName string, routeTablesSummaryName string) RouteTablesSummaryId {
	return RouteTablesSummaryId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		RouteTablesSummaryName:  routeTablesSummaryName,
	}
}

// ParseRouteTablesSummaryID parses 'input' into a RouteTablesSummaryId
func ParseRouteTablesSummaryID(input string) (*RouteTablesSummaryId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteTablesSummaryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteTablesSummaryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCircuitName, ok = parsed.Parsed["expressRouteCircuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	if id.RouteTablesSummaryName, ok = parsed.Parsed["routeTablesSummaryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTablesSummaryName", *parsed)
	}

	return &id, nil
}

// ParseRouteTablesSummaryIDInsensitively parses 'input' case-insensitively into a RouteTablesSummaryId
// note: this method should only be used for API response data and not user input
func ParseRouteTablesSummaryIDInsensitively(input string) (*RouteTablesSummaryId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteTablesSummaryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteTablesSummaryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCircuitName, ok = parsed.Parsed["expressRouteCircuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	if id.RouteTablesSummaryName, ok = parsed.Parsed["routeTablesSummaryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTablesSummaryName", *parsed)
	}

	return &id, nil
}

// ValidateRouteTablesSummaryID checks that 'input' can be parsed as a Route Tables Summary ID
func ValidateRouteTablesSummaryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteTablesSummaryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Tables Summary ID
func (id RouteTablesSummaryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/routeTablesSummary/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName, id.RouteTablesSummaryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Tables Summary ID
func (id RouteTablesSummaryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCircuits", "expressRouteCircuits", "expressRouteCircuits"),
		resourceids.UserSpecifiedSegment("expressRouteCircuitName", "expressRouteCircuitValue"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringValue"),
		resourceids.StaticSegment("staticRouteTablesSummary", "routeTablesSummary", "routeTablesSummary"),
		resourceids.UserSpecifiedSegment("routeTablesSummaryName", "routeTablesSummaryValue"),
	}
}

// String returns a human-readable description of this Route Tables Summary ID
func (id RouteTablesSummaryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Route Tables Summary Name: %q", id.RouteTablesSummaryName),
	}
	return fmt.Sprintf("Route Tables Summary (%s)", strings.Join(components, "\n"))
}
