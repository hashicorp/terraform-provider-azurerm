package expressrouteproviderports

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRouteProviderPortId{})
}

var _ resourceids.ResourceId = &ExpressRouteProviderPortId{}

// ExpressRouteProviderPortId is a struct representing the Resource ID for a Express Route Provider Port
type ExpressRouteProviderPortId struct {
	SubscriptionId               string
	ExpressRouteProviderPortName string
}

// NewExpressRouteProviderPortID returns a new ExpressRouteProviderPortId struct
func NewExpressRouteProviderPortID(subscriptionId string, expressRouteProviderPortName string) ExpressRouteProviderPortId {
	return ExpressRouteProviderPortId{
		SubscriptionId:               subscriptionId,
		ExpressRouteProviderPortName: expressRouteProviderPortName,
	}
}

// ParseExpressRouteProviderPortID parses 'input' into a ExpressRouteProviderPortId
func ParseExpressRouteProviderPortID(input string) (*ExpressRouteProviderPortId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteProviderPortId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteProviderPortId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRouteProviderPortIDInsensitively parses 'input' case-insensitively into a ExpressRouteProviderPortId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteProviderPortIDInsensitively(input string) (*ExpressRouteProviderPortId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteProviderPortId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteProviderPortId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRouteProviderPortId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ExpressRouteProviderPortName, ok = input.Parsed["expressRouteProviderPortName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteProviderPortName", input)
	}

	return nil
}

// ValidateExpressRouteProviderPortID checks that 'input' can be parsed as a Express Route Provider Port ID
func ValidateExpressRouteProviderPortID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteProviderPortID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Provider Port ID
func (id ExpressRouteProviderPortId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/expressRouteProviderPorts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ExpressRouteProviderPortName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Provider Port ID
func (id ExpressRouteProviderPortId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteProviderPorts", "expressRouteProviderPorts", "expressRouteProviderPorts"),
		resourceids.UserSpecifiedSegment("expressRouteProviderPortName", "expressRouteProviderPortName"),
	}
}

// String returns a human-readable description of this Express Route Provider Port ID
func (id ExpressRouteProviderPortId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Express Route Provider Port Name: %q", id.ExpressRouteProviderPortName),
	}
	return fmt.Sprintf("Express Route Provider Port (%s)", strings.Join(components, "\n"))
}
