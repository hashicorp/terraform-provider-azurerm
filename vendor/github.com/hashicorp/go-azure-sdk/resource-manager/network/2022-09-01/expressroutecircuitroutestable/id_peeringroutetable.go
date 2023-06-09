package expressroutecircuitroutestable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PeeringRouteTableId{}

// PeeringRouteTableId is a struct representing the Resource ID for a Peering Route Table
type PeeringRouteTableId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
	PeeringName             string
	RouteTableName          string
}

// NewPeeringRouteTableID returns a new PeeringRouteTableId struct
func NewPeeringRouteTableID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string, peeringName string, routeTableName string) PeeringRouteTableId {
	return PeeringRouteTableId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		RouteTableName:          routeTableName,
	}
}

// ParsePeeringRouteTableID parses 'input' into a PeeringRouteTableId
func ParsePeeringRouteTableID(input string) (*PeeringRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(PeeringRouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PeeringRouteTableId{}

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

	if id.RouteTableName, ok = parsed.Parsed["routeTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", *parsed)
	}

	return &id, nil
}

// ParsePeeringRouteTableIDInsensitively parses 'input' case-insensitively into a PeeringRouteTableId
// note: this method should only be used for API response data and not user input
func ParsePeeringRouteTableIDInsensitively(input string) (*PeeringRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(PeeringRouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PeeringRouteTableId{}

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

	if id.RouteTableName, ok = parsed.Parsed["routeTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", *parsed)
	}

	return &id, nil
}

// ValidatePeeringRouteTableID checks that 'input' can be parsed as a Peering Route Table ID
func ValidatePeeringRouteTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeeringRouteTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Peering Route Table ID
func (id PeeringRouteTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/routeTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName, id.RouteTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Peering Route Table ID
func (id PeeringRouteTableId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableValue"),
	}
}

// String returns a human-readable description of this Peering Route Table ID
func (id PeeringRouteTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Route Table Name: %q", id.RouteTableName),
	}
	return fmt.Sprintf("Peering Route Table (%s)", strings.Join(components, "\n"))
}
