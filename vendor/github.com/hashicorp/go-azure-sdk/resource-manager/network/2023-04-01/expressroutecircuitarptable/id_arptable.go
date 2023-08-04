package expressroutecircuitarptable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ArpTableId{}

// ArpTableId is a struct representing the Resource ID for a Arp Table
type ArpTableId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
	PeeringName             string
	ArpTableName            string
}

// NewArpTableID returns a new ArpTableId struct
func NewArpTableID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string, peeringName string, arpTableName string) ArpTableId {
	return ArpTableId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		ArpTableName:            arpTableName,
	}
}

// ParseArpTableID parses 'input' into a ArpTableId
func ParseArpTableID(input string) (*ArpTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(ArpTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ArpTableId{}

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

	if id.ArpTableName, ok = parsed.Parsed["arpTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "arpTableName", *parsed)
	}

	return &id, nil
}

// ParseArpTableIDInsensitively parses 'input' case-insensitively into a ArpTableId
// note: this method should only be used for API response data and not user input
func ParseArpTableIDInsensitively(input string) (*ArpTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(ArpTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ArpTableId{}

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

	if id.ArpTableName, ok = parsed.Parsed["arpTableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "arpTableName", *parsed)
	}

	return &id, nil
}

// ValidateArpTableID checks that 'input' can be parsed as a Arp Table ID
func ValidateArpTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseArpTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Arp Table ID
func (id ArpTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/arpTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName, id.ArpTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Arp Table ID
func (id ArpTableId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticArpTables", "arpTables", "arpTables"),
		resourceids.UserSpecifiedSegment("arpTableName", "arpTableValue"),
	}
}

// String returns a human-readable description of this Arp Table ID
func (id ArpTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Arp Table Name: %q", id.ArpTableName),
	}
	return fmt.Sprintf("Arp Table (%s)", strings.Join(components, "\n"))
}
