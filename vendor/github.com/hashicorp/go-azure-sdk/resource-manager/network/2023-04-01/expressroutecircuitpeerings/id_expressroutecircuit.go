package expressroutecircuitpeerings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExpressRouteCircuitId{}

// ExpressRouteCircuitId is a struct representing the Resource ID for a Express Route Circuit
type ExpressRouteCircuitId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
}

// NewExpressRouteCircuitID returns a new ExpressRouteCircuitId struct
func NewExpressRouteCircuitID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string) ExpressRouteCircuitId {
	return ExpressRouteCircuitId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
	}
}

// ParseExpressRouteCircuitID parses 'input' into a ExpressRouteCircuitId
func ParseExpressRouteCircuitID(input string) (*ExpressRouteCircuitId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCircuitId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCircuitId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCircuitName, ok = parsed.Parsed["expressRouteCircuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", *parsed)
	}

	return &id, nil
}

// ParseExpressRouteCircuitIDInsensitively parses 'input' case-insensitively into a ExpressRouteCircuitId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteCircuitIDInsensitively(input string) (*ExpressRouteCircuitId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCircuitId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCircuitId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCircuitName, ok = parsed.Parsed["expressRouteCircuitName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRouteCircuitID checks that 'input' can be parsed as a Express Route Circuit ID
func ValidateExpressRouteCircuitID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteCircuitID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Circuit ID
func (id ExpressRouteCircuitId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Circuit ID
func (id ExpressRouteCircuitId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCircuits", "expressRouteCircuits", "expressRouteCircuits"),
		resourceids.UserSpecifiedSegment("expressRouteCircuitName", "expressRouteCircuitValue"),
	}
}

// String returns a human-readable description of this Express Route Circuit ID
func (id ExpressRouteCircuitId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
	}
	return fmt.Sprintf("Express Route Circuit (%s)", strings.Join(components, "\n"))
}
