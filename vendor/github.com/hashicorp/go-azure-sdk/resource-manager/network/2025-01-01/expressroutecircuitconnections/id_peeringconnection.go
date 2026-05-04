package expressroutecircuitconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PeeringConnectionId{})
}

var _ resourceids.ResourceId = &PeeringConnectionId{}

// PeeringConnectionId is a struct representing the Resource ID for a Peering Connection
type PeeringConnectionId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
	PeeringName             string
	ConnectionName          string
}

// NewPeeringConnectionID returns a new PeeringConnectionId struct
func NewPeeringConnectionID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string, peeringName string, connectionName string) PeeringConnectionId {
	return PeeringConnectionId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		ConnectionName:          connectionName,
	}
}

// ParsePeeringConnectionID parses 'input' into a PeeringConnectionId
func ParsePeeringConnectionID(input string) (*PeeringConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeeringConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePeeringConnectionIDInsensitively parses 'input' case-insensitively into a PeeringConnectionId
// note: this method should only be used for API response data and not user input
func ParsePeeringConnectionIDInsensitively(input string) (*PeeringConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeeringConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeeringConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PeeringConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConnectionName, ok = input.Parsed["connectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionName", input)
	}

	return nil
}

// ValidatePeeringConnectionID checks that 'input' can be parsed as a Peering Connection ID
func ValidatePeeringConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeeringConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Peering Connection ID
func (id PeeringConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Peering Connection ID
func (id PeeringConnectionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticConnections", "connections", "connections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionName"),
	}
}

// String returns a human-readable description of this Peering Connection ID
func (id PeeringConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("Peering Connection (%s)", strings.Join(components, "\n"))
}
