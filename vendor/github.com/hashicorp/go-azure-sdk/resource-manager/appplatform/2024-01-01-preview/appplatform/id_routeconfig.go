package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteConfigId{})
}

var _ resourceids.ResourceId = &RouteConfigId{}

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
	parser := resourceids.NewParserFromResourceIdType(&RouteConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteConfigIDInsensitively parses 'input' case-insensitively into a RouteConfigId
// note: this method should only be used for API response data and not user input
func ParseRouteConfigIDInsensitively(input string) (*RouteConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteConfigId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.GatewayName, ok = input.Parsed["gatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", input)
	}

	if id.RouteConfigName, ok = input.Parsed["routeConfigName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeConfigName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayName"),
		resourceids.StaticSegment("staticRouteConfigs", "routeConfigs", "routeConfigs"),
		resourceids.UserSpecifiedSegment("routeConfigName", "routeConfigName"),
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
