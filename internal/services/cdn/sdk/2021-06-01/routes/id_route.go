package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RouteId{}

// RouteId is a struct representing the Resource ID for a Route
type RouteId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
	RouteName         string
}

// NewRouteID returns a new RouteId struct
func NewRouteID(subscriptionId string, resourceGroupName string, profileName string, endpointName string, routeName string) RouteId {
	return RouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
		RouteName:         routeName,
	}
}

// ParseRouteID parses 'input' into a RouteId
func ParseRouteID(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.EndpointName, ok = parsed.Parsed["endpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'endpointName' was not found in the resource id %q", input)
	}

	if id.RouteName, ok = parsed.Parsed["routeName"]; !ok {
		return nil, fmt.Errorf("the segment 'routeName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRouteIDInsensitively parses 'input' case-insensitively into a RouteId
// note: this method should only be used for API response data and not user input
func ParseRouteIDInsensitively(input string) (*RouteId, error) {
	parser := resourceids.NewParserFromResourceIdType(RouteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RouteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.EndpointName, ok = parsed.Parsed["endpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'endpointName' was not found in the resource id %q", input)
	}

	if id.RouteName, ok = parsed.Parsed["routeName"]; !ok {
		return nil, fmt.Errorf("the segment 'routeName' was not found in the resource id %q", input)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/afdEndpoints/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.RouteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Route ID
func (id RouteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticAfdEndpoints", "afdEndpoints", "afdEndpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointValue"),
		resourceids.StaticSegment("staticRoutes", "routes", "routes"),
		resourceids.UserSpecifiedSegment("routeName", "routeValue"),
	}
}

// String returns a human-readable description of this Route ID
func (id RouteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
		fmt.Sprintf("Route Name: %q", id.RouteName),
	}
	return fmt.Sprintf("Route (%s)", strings.Join(components, "\n"))
}
