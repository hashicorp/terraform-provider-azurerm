package expressroutecrossconnectionroutetable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRouteCrossConnectionPeeringRouteTableId{})
}

var _ resourceids.ResourceId = &ExpressRouteCrossConnectionPeeringRouteTableId{}

// ExpressRouteCrossConnectionPeeringRouteTableId is a struct representing the Resource ID for a Express Route Cross Connection Peering Route Table
type ExpressRouteCrossConnectionPeeringRouteTableId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	ExpressRouteCrossConnectionName string
	PeeringName                     string
	RouteTableName                  string
}

// NewExpressRouteCrossConnectionPeeringRouteTableID returns a new ExpressRouteCrossConnectionPeeringRouteTableId struct
func NewExpressRouteCrossConnectionPeeringRouteTableID(subscriptionId string, resourceGroupName string, expressRouteCrossConnectionName string, peeringName string, routeTableName string) ExpressRouteCrossConnectionPeeringRouteTableId {
	return ExpressRouteCrossConnectionPeeringRouteTableId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		ExpressRouteCrossConnectionName: expressRouteCrossConnectionName,
		PeeringName:                     peeringName,
		RouteTableName:                  routeTableName,
	}
}

// ParseExpressRouteCrossConnectionPeeringRouteTableID parses 'input' into a ExpressRouteCrossConnectionPeeringRouteTableId
func ParseExpressRouteCrossConnectionPeeringRouteTableID(input string) (*ExpressRouteCrossConnectionPeeringRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteCrossConnectionPeeringRouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteCrossConnectionPeeringRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRouteCrossConnectionPeeringRouteTableIDInsensitively parses 'input' case-insensitively into a ExpressRouteCrossConnectionPeeringRouteTableId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteCrossConnectionPeeringRouteTableIDInsensitively(input string) (*ExpressRouteCrossConnectionPeeringRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteCrossConnectionPeeringRouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteCrossConnectionPeeringRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRouteCrossConnectionPeeringRouteTableId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteCrossConnectionName, ok = input.Parsed["expressRouteCrossConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", input)
	}

	if id.PeeringName, ok = input.Parsed["peeringName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "peeringName", input)
	}

	if id.RouteTableName, ok = input.Parsed["routeTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", input)
	}

	return nil
}

// ValidateExpressRouteCrossConnectionPeeringRouteTableID checks that 'input' can be parsed as a Express Route Cross Connection Peering Route Table ID
func ValidateExpressRouteCrossConnectionPeeringRouteTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteCrossConnectionPeeringRouteTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Cross Connection Peering Route Table ID
func (id ExpressRouteCrossConnectionPeeringRouteTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCrossConnections/%s/peerings/%s/routeTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCrossConnectionName, id.PeeringName, id.RouteTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Cross Connection Peering Route Table ID
func (id ExpressRouteCrossConnectionPeeringRouteTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCrossConnections", "expressRouteCrossConnections", "expressRouteCrossConnections"),
		resourceids.UserSpecifiedSegment("expressRouteCrossConnectionName", "expressRouteCrossConnectionName"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringName"),
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableName"),
	}
}

// String returns a human-readable description of this Express Route Cross Connection Peering Route Table ID
func (id ExpressRouteCrossConnectionPeeringRouteTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Cross Connection Name: %q", id.ExpressRouteCrossConnectionName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Route Table Name: %q", id.RouteTableName),
	}
	return fmt.Sprintf("Express Route Cross Connection Peering Route Table (%s)", strings.Join(components, "\n"))
}
