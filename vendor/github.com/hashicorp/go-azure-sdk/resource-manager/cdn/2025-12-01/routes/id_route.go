package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RouteId{})
}

var _ resourceids.ResourceId = &RouteId{}

// RouteId is a struct representing the Resource ID for a Route
type RouteId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	AfdEndpointName   string
	RouteName         string
}

// NewRouteID returns a new RouteId struct
func NewRouteID(subscriptionId string, resourceGroupName string, profileName string, afdEndpointName string, routeName string) RouteId {
	return RouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		AfdEndpointName:   afdEndpointName,
		RouteName:         routeName,
	}
}

// ParseRouteID parses 'input' into a RouteId
func ParseRouteID(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRouteIDInsensitively parses 'input' case-insensitively into a RouteId
// note: this method should only be used for API response data and not user input
func ParseRouteIDInsensitively(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RouteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RouteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RouteId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.AfdEndpointName, ok = input.Parsed["afdEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "afdEndpointName", input)
	}

	if id.RouteName, ok = input.Parsed["routeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "routeName", input)
	}

	return nil
}

// ValidateRouteID checks that 'input' can be parsed as a Route ID
func ValidateRouteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRouteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Route ID
func (id RouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.AfdEndpointName, id.RouteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route ID
func (id RouteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticAfdEndpoints", "afdEndpoints", "afdEndpoints"),
		resourceids.UserSpecifiedSegment("afdEndpointName", "afdEndpointName"),
		resourceids.StaticSegment("staticRoutes", "routes", "routes"),
		resourceids.UserSpecifiedSegment("routeName", "routeName"),
	}
}

// String returns a human-readable description of this Route ID
func (id RouteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Afd Endpoint Name: %q", id.AfdEndpointName),
		fmt.Sprintf("Route Name: %q", id.RouteName),
	}
	return fmt.Sprintf("Route (%s)", strings.Join(components, "\n"))
}
