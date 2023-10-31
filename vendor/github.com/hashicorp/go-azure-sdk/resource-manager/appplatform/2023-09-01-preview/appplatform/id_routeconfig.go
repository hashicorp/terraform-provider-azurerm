package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RouteConfigId{}

// RouteConfigId is a struct representing the Resource ID for a Route Config
type RouteConfigId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	GatewayName       string
	RouteConfigName   string
}

// NewRouteConfigID returns a new RouteConfigId struct
func NewRouteConfigID(subscriptionId string, resourceGroupName string, springName string, gatewayName string, routeConfigName string) RouteConfigId {
	return RouteConfigId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		GatewayName:       gatewayName,
		RouteConfigName:   routeConfigName,
	}
}

// ParseRouteConfigID parses 'input' into a RouteConfigId
func ParseRouteConfigID(input string) (*RouteConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", *parsed)
	}

	if id.RouteConfigName, ok = parsed.Parsed["routeConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeConfigName", *parsed)
	}

	return &id, nil
}

// ParseRouteConfigIDInsensitively parses 'input' case-insensitively into a RouteConfigId
// note: this method should only be used for API response data and not user input
func ParseRouteConfigIDInsensitively(input string) (*RouteConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", *parsed)
	}

	if id.RouteConfigName, ok = parsed.Parsed["routeConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "routeConfigName", *parsed)
	}

	return &id, nil
}

// ValidateRouteConfigID checks that 'input' can be parsed as a Route Config ID
func ValidateRouteConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route Config ID
func (id RouteConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/gateways/%s/routeConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route Config ID
func (id RouteConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayValue"),
		resourceids.StaticSegment("staticRouteConfigs", "routeConfigs", "routeConfigs"),
		resourceids.UserSpecifiedSegment("routeConfigName", "routeConfigValue"),
	}
}

// String returns a human-readable description of this Route Config ID
func (id RouteConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
		fmt.Sprintf("Route Config Name: %q", id.RouteConfigName),
	}
	return fmt.Sprintf("Route Config (%s)", strings.Join(components, "\n"))
}
