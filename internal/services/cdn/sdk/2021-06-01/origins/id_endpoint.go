package origins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EndpointId{}

// EndpointId is a struct representing the Resource ID for a Endpoint
type EndpointId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
}

// NewEndpointID returns a new EndpointId struct
func NewEndpointID(subscriptionId string, resourceGroupName string, profileName string, endpointName string) EndpointId {
	return EndpointId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
	}
}

// ParseEndpointID parses 'input' into a EndpointId
func ParseEndpointID(input string) (*EndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointId{}

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

	return &id, nil
}

// ParseEndpointIDInsensitively parses 'input' case-insensitively into a EndpointId
// note: this method should only be used for API response data and not user input
func ParseEndpointIDInsensitively(input string) (*EndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointId{}

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

	return &id, nil
}

// ValidateEndpointID checks that 'input' can be parsed as a Endpoint ID
func ValidateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Endpoint ID
func (id EndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Endpoint ID
func (id EndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticEndpoints", "endpoints", "endpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointValue"),
	}
}

// String returns a human-readable description of this Endpoint ID
func (id EndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
	}
	return fmt.Sprintf("Endpoint (%s)", strings.Join(components, "\n"))
}
