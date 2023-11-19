package expressroutecrossconnectionroutetablesummary

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PeeringRouteTablesSummaryId{}

// PeeringRouteTablesSummaryId is a struct representing the Resource ID for a Peering Route Tables Summary
type PeeringRouteTablesSummaryId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	ExpressRouteCrossConnectionName string
	PeeringName                     string
	RouteTablesSummaryName          string
}

// NewPeeringRouteTablesSummaryID returns a new PeeringRouteTablesSummaryId struct
func NewPeeringRouteTablesSummaryID(subscriptionId string, resourceGroupName string, expressRouteCrossConnectionName string, peeringName string, routeTablesSummaryName string) PeeringRouteTablesSummaryId {
	return PeeringRouteTablesSummaryId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		ExpressRouteCrossConnectionName: expressRouteCrossConnectionName,
		PeeringName:                     peeringName,
		RouteTablesSummaryName:          routeTablesSummaryName,
	}
}

// ParsePeeringRouteTablesSummaryID parses 'input' into a PeeringRouteTablesSummaryId
func ParsePeeringRouteTablesSummaryID(input string) (*PeeringRouteTablesSummaryId, error) {
	parser := resourceids.NewParserFromResourceIdType(PeeringRouteTablesSummaryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PeeringRouteTablesSummaryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCrossConnectionName, ok = parsed.Parsed["expressRouteCrossConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	if id.RouteTablesSummaryName, ok = parsed.Parsed["routeTablesSummaryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTablesSummaryName", *parsed)
	}

	return &id, nil
}

// ParsePeeringRouteTablesSummaryIDInsensitively parses 'input' case-insensitively into a PeeringRouteTablesSummaryId
// note: this method should only be used for API response data and not user input
func ParsePeeringRouteTablesSummaryIDInsensitively(input string) (*PeeringRouteTablesSummaryId, error) {
	parser := resourceids.NewParserFromResourceIdType(PeeringRouteTablesSummaryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PeeringRouteTablesSummaryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCrossConnectionName, ok = parsed.Parsed["expressRouteCrossConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	if id.RouteTablesSummaryName, ok = parsed.Parsed["routeTablesSummaryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeTablesSummaryName", *parsed)
	}

	return &id, nil
}

// ValidatePeeringRouteTablesSummaryID checks that 'input' can be parsed as a Peering Route Tables Summary ID
func ValidatePeeringRouteTablesSummaryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeeringRouteTablesSummaryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Peering Route Tables Summary ID
func (id PeeringRouteTablesSummaryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCrossConnections/%s/peerings/%s/routeTablesSummary/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCrossConnectionName, id.PeeringName, id.RouteTablesSummaryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Peering Route Tables Summary ID
func (id PeeringRouteTablesSummaryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCrossConnections", "expressRouteCrossConnections", "expressRouteCrossConnections"),
		resourceids.UserSpecifiedSegment("expressRouteCrossConnectionName", "expressRouteCrossConnectionValue"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringValue"),
		resourceids.StaticSegment("staticRouteTablesSummary", "routeTablesSummary", "routeTablesSummary"),
		resourceids.UserSpecifiedSegment("routeTablesSummaryName", "routeTablesSummaryValue"),
	}
}

// String returns a human-readable description of this Peering Route Tables Summary ID
func (id PeeringRouteTablesSummaryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Cross Connection Name: %q", id.ExpressRouteCrossConnectionName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Route Tables Summary Name: %q", id.RouteTablesSummaryName),
	}
	return fmt.Sprintf("Peering Route Tables Summary (%s)", strings.Join(components, "\n"))
}
