package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualHubRouteTableId{})
}

var _ resourceids.ResourceId = &VirtualHubRouteTableId{}

// VirtualHubRouteTableId is a struct representing the Resource ID for a Virtual Hub Route Table
type VirtualHubRouteTableId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualHubName    string
	RouteTableName    string
}

// NewVirtualHubRouteTableID returns a new VirtualHubRouteTableId struct
func NewVirtualHubRouteTableID(subscriptionId string, resourceGroupName string, virtualHubName string, routeTableName string) VirtualHubRouteTableId {
	return VirtualHubRouteTableId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualHubName:    virtualHubName,
		RouteTableName:    routeTableName,
	}
}

// ParseVirtualHubRouteTableID parses 'input' into a VirtualHubRouteTableId
func ParseVirtualHubRouteTableID(input string) (*VirtualHubRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualHubRouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualHubRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualHubRouteTableIDInsensitively parses 'input' case-insensitively into a VirtualHubRouteTableId
// note: this method should only be used for API response data and not user input
func ParseVirtualHubRouteTableIDInsensitively(input string) (*VirtualHubRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualHubRouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualHubRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualHubRouteTableId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualHubName, ok = input.Parsed["virtualHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", input)
	}

	if id.RouteTableName, ok = input.Parsed["routeTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeTableName", input)
	}

	return nil
}

// ValidateVirtualHubRouteTableID checks that 'input' can be parsed as a Virtual Hub Route Table ID
func ValidateVirtualHubRouteTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualHubRouteTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Hub Route Table ID
func (id VirtualHubRouteTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/routeTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.RouteTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Hub Route Table ID
func (id VirtualHubRouteTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticRouteTables", "routeTables", "routeTables"),
		resourceids.UserSpecifiedSegment("routeTableName", "routeTableName"),
	}
}

// String returns a human-readable description of this Virtual Hub Route Table ID
func (id VirtualHubRouteTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Route Table Name: %q", id.RouteTableName),
	}
	return fmt.Sprintf("Virtual Hub Route Table (%s)", strings.Join(components, "\n"))
}
