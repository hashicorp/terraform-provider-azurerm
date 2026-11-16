package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteTableId{})
}

var _ resourceids.ResourceId = &RouteTableId{}

// RouteTableId is a struct representing the Resource ID for a Route Table
type RouteTableId struct {
	SubscriptionId    string
	ResourceGroupName string
	RouteTableName    string
}

// NewRouteTableID returns a new RouteTableId struct
func NewRouteTableID(subscriptionId string, resourceGroupName string, routeTableName string) RouteTableId {
	return RouteTableId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RouteTableName:    routeTableName,
	}
}

// ParseRouteTableID parses 'input' into a RouteTableId
func ParseRouteTableID(input string) (*RouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteTableIDInsensitively parses 'input' case-insensitively into a RouteTableId
// note: this method should only be used for API response data and not user input
func ParseRouteTableIDInsensitively(input string) (*RouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteTableId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RouteTableName, ok = input.Parsed["routeTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", input)
	}

	return nil
}

// ValidateRouteTableID checks that 'input' can be parsed as a Route Table ID
func ValidateRouteTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Table ID
func (id RouteTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RouteTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Table ID
func (id RouteTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableName"),
	}
}

// String returns a human-readable description of this Route Table ID
func (id RouteTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Route Table Name: %q", id.RouteTableName),
	}
	return fmt.Sprintf("Route Table (%s)", strings.Join(components, "\n"))
}
