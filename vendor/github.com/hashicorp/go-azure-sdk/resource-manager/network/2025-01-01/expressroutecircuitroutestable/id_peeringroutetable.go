package expressroutecircuitroutestable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PeeringRouteTableId{})
}

var _ resourceids.ResourceId = &PeeringRouteTableId{}

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
	parser := resourceids.NewParserFromResourceIdType(&PeeringRouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePeeringRouteTableIDInsensitively parses 'input' case-insensitively into a PeeringRouteTableId
// note: this method should only be used for API response data and not user input
func ParsePeeringRouteTableIDInsensitively(input string) (*PeeringRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeeringRouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PeeringRouteTableId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteCircuitName, ok = input.Parsed["expressRouteCircuitName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", input)
	}

	if id.PeeringName, ok = input.Parsed["peeringName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "peeringName", input)
	}

	if id.RouteTableName, ok = input.Parsed["routeTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("expressRouteCircuitName", "expressRouteCircuitName"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringName"),
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableName"),
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
