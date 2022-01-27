package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AfdEndpointId{}

// AfdEndpointId is a struct representing the Resource ID for a Afd Endpoint
type AfdEndpointId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
}

// NewAfdEndpointID returns a new AfdEndpointId struct
func NewAfdEndpointID(subscriptionId string, resourceGroupName string, profileName string, endpointName string) AfdEndpointId {
	return AfdEndpointId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
	}
}

// ParseAfdEndpointID parses 'input' into a AfdEndpointId
func ParseAfdEndpointID(input string) (*AfdEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(AfdEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AfdEndpointId{}

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

// ParseAfdEndpointIDInsensitively parses 'input' case-insensitively into a AfdEndpointId
// note: this method should only be used for API response data and not user input
func ParseAfdEndpointIDInsensitively(input string) (*AfdEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(AfdEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AfdEndpointId{}

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

// ValidateAfdEndpointID checks that 'input' can be parsed as a Afd Endpoint ID
func ValidateAfdEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAfdEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Afd Endpoint ID
func (id AfdEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/afdEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Afd Endpoint ID
func (id AfdEndpointId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Afd Endpoint ID
func (id AfdEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
	}
	return fmt.Sprintf("Afd Endpoint (%s)", strings.Join(components, "\n"))
}
