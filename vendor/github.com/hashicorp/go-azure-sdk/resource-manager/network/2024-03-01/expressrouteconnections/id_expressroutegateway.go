package expressrouteconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRouteGatewayId{})
}

var _ resourceids.ResourceId = &ExpressRouteGatewayId{}

// ExpressRouteGatewayId is a struct representing the Resource ID for a Express Route Gateway
type ExpressRouteGatewayId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteGatewayName string
}

// NewExpressRouteGatewayID returns a new ExpressRouteGatewayId struct
func NewExpressRouteGatewayID(subscriptionId string, resourceGroupName string, expressRouteGatewayName string) ExpressRouteGatewayId {
	return ExpressRouteGatewayId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteGatewayName: expressRouteGatewayName,
	}
}

// ParseExpressRouteGatewayID parses 'input' into a ExpressRouteGatewayId
func ParseExpressRouteGatewayID(input string) (*ExpressRouteGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRouteGatewayIDInsensitively parses 'input' case-insensitively into a ExpressRouteGatewayId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteGatewayIDInsensitively(input string) (*ExpressRouteGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRouteGatewayId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteGatewayName, ok = input.Parsed["expressRouteGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteGatewayName", input)
	}

	return nil
}

// ValidateExpressRouteGatewayID checks that 'input' can be parsed as a Express Route Gateway ID
func ValidateExpressRouteGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Gateway ID
func (id ExpressRouteGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Gateway ID
func (id ExpressRouteGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteGateways", "expressRouteGateways", "expressRouteGateways"),
		resourceids.UserSpecifiedSegment("expressRouteGatewayName", "expressRouteGatewayName"),
	}
}

// String returns a human-readable description of this Express Route Gateway ID
func (id ExpressRouteGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Gateway Name: %q", id.ExpressRouteGatewayName),
	}
	return fmt.Sprintf("Express Route Gateway (%s)", strings.Join(components, "\n"))
}
