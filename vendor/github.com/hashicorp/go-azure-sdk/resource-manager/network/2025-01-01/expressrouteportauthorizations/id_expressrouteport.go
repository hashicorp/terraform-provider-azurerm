package expressrouteportauthorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRoutePortId{})
}

var _ resourceids.ResourceId = &ExpressRoutePortId{}

// ExpressRoutePortId is a struct representing the Resource ID for a Express Route Port
type ExpressRoutePortId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ExpressRoutePortName string
}

// NewExpressRoutePortID returns a new ExpressRoutePortId struct
func NewExpressRoutePortID(subscriptionId string, resourceGroupName string, expressRoutePortName string) ExpressRoutePortId {
	return ExpressRoutePortId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ExpressRoutePortName: expressRoutePortName,
	}
}

// ParseExpressRoutePortID parses 'input' into a ExpressRoutePortId
func ParseExpressRoutePortID(input string) (*ExpressRoutePortId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRoutePortId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRoutePortId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRoutePortIDInsensitively parses 'input' case-insensitively into a ExpressRoutePortId
// note: this method should only be used for API response data and not user input
func ParseExpressRoutePortIDInsensitively(input string) (*ExpressRoutePortId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRoutePortId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRoutePortId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRoutePortId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRoutePortName, ok = input.Parsed["expressRoutePortName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortName", input)
	}

	return nil
}

// ValidateExpressRoutePortID checks that 'input' can be parsed as a Express Route Port ID
func ValidateExpressRoutePortID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRoutePortID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Port ID
func (id ExpressRoutePortId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRoutePorts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Port ID
func (id ExpressRoutePortId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRoutePorts", "expressRoutePorts", "expressRoutePorts"),
		resourceids.UserSpecifiedSegment("expressRoutePortName", "expressRoutePortName"),
	}
}

// String returns a human-readable description of this Express Route Port ID
func (id ExpressRoutePortId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Port Name: %q", id.ExpressRoutePortName),
	}
	return fmt.Sprintf("Express Route Port (%s)", strings.Join(components, "\n"))
}
