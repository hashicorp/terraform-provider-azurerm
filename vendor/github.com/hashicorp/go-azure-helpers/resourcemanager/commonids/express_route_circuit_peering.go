// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ExpressRouteCircuitPeeringId{}

// ExpressRouteCircuitPeeringId is a struct representing the Resource ID for a Express Route Circuit Peering
type ExpressRouteCircuitPeeringId struct {
	SubscriptionId    string
	ResourceGroupName string
	CircuitName       string
	PeeringName       string
}

// NewExpressRouteCircuitPeeringID returns a new ExpressRouteCircuitPeeringId struct
func NewExpressRouteCircuitPeeringID(subscriptionId string, resourceGroupName string, circuitName string, peeringName string) ExpressRouteCircuitPeeringId {
	return ExpressRouteCircuitPeeringId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CircuitName:       circuitName,
		PeeringName:       peeringName,
	}
}

// ParseExpressRouteCircuitPeeringID parses 'input' into a ExpressRouteCircuitPeeringId
func ParseExpressRouteCircuitPeeringID(input string) (*ExpressRouteCircuitPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCircuitPeeringId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCircuitPeeringId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CircuitName, ok = parsed.Parsed["circuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "circuitName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ParseExpressRouteCircuitPeeringIDInsensitively parses 'input' case-insensitively into a ExpressRouteCircuitPeeringId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteCircuitPeeringIDInsensitively(input string) (*ExpressRouteCircuitPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCircuitPeeringId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCircuitPeeringId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CircuitName, ok = parsed.Parsed["circuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "circuitName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRouteCircuitPeeringID checks that 'input' can be parsed as a Express Route Circuit Peering ID
func ValidateExpressRouteCircuitPeeringID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteCircuitPeeringID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Circuit Peering ID
func (id ExpressRouteCircuitPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CircuitName, id.PeeringName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Circuit Peering ID
func (id ExpressRouteCircuitPeeringId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("expressRouteCircuits", "expressRouteCircuits", "expressRouteCircuits"),
		resourceids.UserSpecifiedSegment("circuitName", "circuitValue"),
		resourceids.StaticSegment("peerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringValue"),
	}
}

// String returns a human-readable description of this Express Route Circuit Peering ID
func (id ExpressRouteCircuitPeeringId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Circuit Name: %q", id.CircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
	}
	return fmt.Sprintf("Express Route Circuit Peering (%s)", strings.Join(components, "\n"))
}
