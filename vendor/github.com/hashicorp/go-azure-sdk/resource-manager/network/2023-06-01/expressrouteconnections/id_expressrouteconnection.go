package expressrouteconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExpressRouteConnectionId{}

// ExpressRouteConnectionId is a struct representing the Resource ID for a Express Route Connection
type ExpressRouteConnectionId struct {
	SubscriptionId             string
	ResourceGroupName          string
	ExpressRouteGatewayName    string
	ExpressRouteConnectionName string
}

// NewExpressRouteConnectionID returns a new ExpressRouteConnectionId struct
func NewExpressRouteConnectionID(subscriptionId string, resourceGroupName string, expressRouteGatewayName string, expressRouteConnectionName string) ExpressRouteConnectionId {
	return ExpressRouteConnectionId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		ExpressRouteGatewayName:    expressRouteGatewayName,
		ExpressRouteConnectionName: expressRouteConnectionName,
	}
}

// ParseExpressRouteConnectionID parses 'input' into a ExpressRouteConnectionId
func ParseExpressRouteConnectionID(input string) (*ExpressRouteConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteGatewayName, ok = parsed.Parsed["expressRouteGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteGatewayName", *parsed)
	}

	if id.ExpressRouteConnectionName, ok = parsed.Parsed["expressRouteConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteConnectionName", *parsed)
	}

	return &id, nil
}

// ParseExpressRouteConnectionIDInsensitively parses 'input' case-insensitively into a ExpressRouteConnectionId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteConnectionIDInsensitively(input string) (*ExpressRouteConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteGatewayName, ok = parsed.Parsed["expressRouteGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteGatewayName", *parsed)
	}

	if id.ExpressRouteConnectionName, ok = parsed.Parsed["expressRouteConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRouteConnectionID checks that 'input' can be parsed as a Express Route Connection ID
func ValidateExpressRouteConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Connection ID
func (id ExpressRouteConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteGateways/%s/expressRouteConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName, id.ExpressRouteConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Connection ID
func (id ExpressRouteConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteGateways", "expressRouteGateways", "expressRouteGateways"),
		resourceids.UserSpecifiedSegment("expressRouteGatewayName", "expressRouteGatewayValue"),
		resourceids.StaticSegment("staticExpressRouteConnections", "expressRouteConnections", "expressRouteConnections"),
		resourceids.UserSpecifiedSegment("expressRouteConnectionName", "expressRouteConnectionValue"),
	}
}

// String returns a human-readable description of this Express Route Connection ID
func (id ExpressRouteConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Gateway Name: %q", id.ExpressRouteGatewayName),
		fmt.Sprintf("Express Route Connection Name: %q", id.ExpressRouteConnectionName),
	}
	return fmt.Sprintf("Express Route Connection (%s)", strings.Join(components, "\n"))
}
