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
	recaser.RegisterResourceId(&HubRouteTableId{})
}

var _ resourceids.ResourceId = &HubRouteTableId{}

// HubRouteTableId is a struct representing the Resource ID for a Hub Route Table
type HubRouteTableId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualHubName    string
	HubRouteTableName string
}

// NewHubRouteTableID returns a new HubRouteTableId struct
func NewHubRouteTableID(subscriptionId string, resourceGroupName string, virtualHubName string, hubRouteTableName string) HubRouteTableId {
	return HubRouteTableId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualHubName:    virtualHubName,
		HubRouteTableName: hubRouteTableName,
	}
}

// ParseHubRouteTableID parses 'input' into a HubRouteTableId
func ParseHubRouteTableID(input string) (*HubRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HubRouteTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHubRouteTableIDInsensitively parses 'input' case-insensitively into a HubRouteTableId
// note: this method should only be used for API response data and not user input
func ParseHubRouteTableIDInsensitively(input string) (*HubRouteTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HubRouteTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubRouteTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HubRouteTableId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HubRouteTableName, ok = input.Parsed["hubRouteTableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hubRouteTableName", input)
	}

	return nil
}

// ValidateHubRouteTableID checks that 'input' can be parsed as a Hub Route Table ID
func ValidateHubRouteTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHubRouteTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hub Route Table ID
func (id HubRouteTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/hubRouteTables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.HubRouteTableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hub Route Table ID
func (id HubRouteTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticHubRouteTables", "hubRouteTables", "hubRouteTables"),
		resourceids.UserSpecifiedSegment("hubRouteTableName", "hubRouteTableName"),
	}
}

// String returns a human-readable description of this Hub Route Table ID
func (id HubRouteTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Hub Route Table Name: %q", id.HubRouteTableName),
	}
	return fmt.Sprintf("Hub Route Table (%s)", strings.Join(components, "\n"))
}
