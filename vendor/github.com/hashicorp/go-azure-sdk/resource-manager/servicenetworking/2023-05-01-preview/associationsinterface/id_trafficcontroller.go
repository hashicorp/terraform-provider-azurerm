package associationsinterface

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TrafficControllerId{}

// TrafficControllerId is a struct representing the Resource ID for a Traffic Controller
type TrafficControllerId struct {
	SubscriptionId        string
	ResourceGroupName     string
	TrafficControllerName string
}

// NewTrafficControllerID returns a new TrafficControllerId struct
func NewTrafficControllerID(subscriptionId string, resourceGroupName string, trafficControllerName string) TrafficControllerId {
	return TrafficControllerId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		TrafficControllerName: trafficControllerName,
	}
}

// ParseTrafficControllerID parses 'input' into a TrafficControllerId
func ParseTrafficControllerID(input string) (*TrafficControllerId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrafficControllerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrafficControllerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	return &id, nil
}

// ParseTrafficControllerIDInsensitively parses 'input' case-insensitively into a TrafficControllerId
// note: this method should only be used for API response data and not user input
func ParseTrafficControllerIDInsensitively(input string) (*TrafficControllerId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrafficControllerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrafficControllerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	return &id, nil
}

// ValidateTrafficControllerID checks that 'input' can be parsed as a Traffic Controller ID
func ValidateTrafficControllerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrafficControllerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Traffic Controller ID
func (id TrafficControllerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceNetworking/trafficControllers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Traffic Controller ID
func (id TrafficControllerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceNetworking", "Microsoft.ServiceNetworking", "Microsoft.ServiceNetworking"),
		resourceids.StaticSegment("staticTrafficControllers", "trafficControllers", "trafficControllers"),
		resourceids.UserSpecifiedSegment("trafficControllerName", "trafficControllerValue"),
	}
}

// String returns a human-readable description of this Traffic Controller ID
func (id TrafficControllerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Traffic Controller Name: %q", id.TrafficControllerName),
	}
	return fmt.Sprintf("Traffic Controller (%s)", strings.Join(components, "\n"))
}
