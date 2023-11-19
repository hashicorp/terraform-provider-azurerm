package frontendsinterface

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FrontendId{}

// FrontendId is a struct representing the Resource ID for a Frontend
type FrontendId struct {
	SubscriptionId        string
	ResourceGroupName     string
	TrafficControllerName string
	FrontendName          string
}

// NewFrontendID returns a new FrontendId struct
func NewFrontendID(subscriptionId string, resourceGroupName string, trafficControllerName string, frontendName string) FrontendId {
	return FrontendId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		TrafficControllerName: trafficControllerName,
		FrontendName:          frontendName,
	}
}

// ParseFrontendID parses 'input' into a FrontendId
func ParseFrontendID(input string) (*FrontendId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	if id.FrontendName, ok = parsed.Parsed["frontendName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendName", *parsed)
	}

	return &id, nil
}

// ParseFrontendIDInsensitively parses 'input' case-insensitively into a FrontendId
// note: this method should only be used for API response data and not user input
func ParseFrontendIDInsensitively(input string) (*FrontendId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	if id.FrontendName, ok = parsed.Parsed["frontendName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendName", *parsed)
	}

	return &id, nil
}

// ValidateFrontendID checks that 'input' can be parsed as a Frontend ID
func ValidateFrontendID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFrontendID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Frontend ID
func (id FrontendId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceNetworking/trafficControllers/%s/frontends/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName, id.FrontendName)
}

// Segments returns a slice of Resource ID Segments which comprise this Frontend ID
func (id FrontendId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceNetworking", "Microsoft.ServiceNetworking", "Microsoft.ServiceNetworking"),
		resourceids.StaticSegment("staticTrafficControllers", "trafficControllers", "trafficControllers"),
		resourceids.UserSpecifiedSegment("trafficControllerName", "trafficControllerValue"),
		resourceids.StaticSegment("staticFrontends", "frontends", "frontends"),
		resourceids.UserSpecifiedSegment("frontendName", "frontendValue"),
	}
}

// String returns a human-readable description of this Frontend ID
func (id FrontendId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Traffic Controller Name: %q", id.TrafficControllerName),
		fmt.Sprintf("Frontend Name: %q", id.FrontendName),
	}
	return fmt.Sprintf("Frontend (%s)", strings.Join(components, "\n"))
}
