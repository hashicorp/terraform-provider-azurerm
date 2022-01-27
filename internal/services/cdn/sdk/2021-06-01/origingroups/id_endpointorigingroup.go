package origingroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EndpointOriginGroupId{}

// EndpointOriginGroupId is a struct representing the Resource ID for a Endpoint Origin Group
type EndpointOriginGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
	OriginGroupName   string
}

// NewEndpointOriginGroupID returns a new EndpointOriginGroupId struct
func NewEndpointOriginGroupID(subscriptionId string, resourceGroupName string, profileName string, endpointName string, originGroupName string) EndpointOriginGroupId {
	return EndpointOriginGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
		OriginGroupName:   originGroupName,
	}
}

// ParseEndpointOriginGroupID parses 'input' into a EndpointOriginGroupId
func ParseEndpointOriginGroupID(input string) (*EndpointOriginGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointOriginGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointOriginGroupId{}

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

	if id.OriginGroupName, ok = parsed.Parsed["originGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'originGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseEndpointOriginGroupIDInsensitively parses 'input' case-insensitively into a EndpointOriginGroupId
// note: this method should only be used for API response data and not user input
func ParseEndpointOriginGroupIDInsensitively(input string) (*EndpointOriginGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointOriginGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointOriginGroupId{}

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

	if id.OriginGroupName, ok = parsed.Parsed["originGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'originGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateEndpointOriginGroupID checks that 'input' can be parsed as a Endpoint Origin Group ID
func ValidateEndpointOriginGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEndpointOriginGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Endpoint Origin Group ID
func (id EndpointOriginGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/endpoints/%s/originGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.OriginGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Endpoint Origin Group ID
func (id EndpointOriginGroupId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticOriginGroups", "originGroups", "originGroups"),
		resourceids.UserSpecifiedSegment("originGroupName", "originGroupValue"),
	}
}

// String returns a human-readable description of this Endpoint Origin Group ID
func (id EndpointOriginGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
		fmt.Sprintf("Origin Group Name: %q", id.OriginGroupName),
	}
	return fmt.Sprintf("Endpoint Origin Group (%s)", strings.Join(components, "\n"))
}
