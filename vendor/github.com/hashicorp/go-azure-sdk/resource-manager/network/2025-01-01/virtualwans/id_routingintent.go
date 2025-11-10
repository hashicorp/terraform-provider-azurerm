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
	recaser.RegisterResourceId(&RoutingIntentId{})
}

var _ resourceids.ResourceId = &RoutingIntentId{}

// RoutingIntentId is a struct representing the Resource ID for a Routing Intent
type RoutingIntentId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualHubName    string
	RoutingIntentName string
}

// NewRoutingIntentID returns a new RoutingIntentId struct
func NewRoutingIntentID(subscriptionId string, resourceGroupName string, virtualHubName string, routingIntentName string) RoutingIntentId {
	return RoutingIntentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualHubName:    virtualHubName,
		RoutingIntentName: routingIntentName,
	}
}

// ParseRoutingIntentID parses 'input' into a RoutingIntentId
func ParseRoutingIntentID(input string) (*RoutingIntentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoutingIntentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoutingIntentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRoutingIntentIDInsensitively parses 'input' case-insensitively into a RoutingIntentId
// note: this method should only be used for API response data and not user input
func ParseRoutingIntentIDInsensitively(input string) (*RoutingIntentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoutingIntentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoutingIntentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RoutingIntentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RoutingIntentName, ok = input.Parsed["routingIntentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routingIntentName", input)
	}

	return nil
}

// ValidateRoutingIntentID checks that 'input' can be parsed as a Routing Intent ID
func ValidateRoutingIntentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoutingIntentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Routing Intent ID
func (id RoutingIntentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/routingIntent/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.RoutingIntentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Routing Intent ID
func (id RoutingIntentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticRoutingIntent", "routingIntent", "routingIntent"),
		resourceids.UserSpecifiedSegment("routingIntentName", "routingIntentName"),
	}
}

// String returns a human-readable description of this Routing Intent ID
func (id RoutingIntentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Routing Intent Name: %q", id.RoutingIntentName),
	}
	return fmt.Sprintf("Routing Intent (%s)", strings.Join(components, "\n"))
}
