package origins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OriginId{}

// OriginId is a struct representing the Resource ID for a Origin
type OriginId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
	OriginName        string
}

// NewOriginID returns a new OriginId struct
func NewOriginID(subscriptionId string, resourceGroupName string, profileName string, endpointName string, originName string) OriginId {
	return OriginId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
		OriginName:        originName,
	}
}

// ParseOriginID parses 'input' into a OriginId
func ParseOriginID(input string) (*OriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(OriginId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OriginId{}

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

	if id.OriginName, ok = parsed.Parsed["originName"]; !ok {
		return nil, fmt.Errorf("the segment 'originName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOriginIDInsensitively parses 'input' case-insensitively into a OriginId
// note: this method should only be used for API response data and not user input
func ParseOriginIDInsensitively(input string) (*OriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(OriginId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OriginId{}

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

	if id.OriginName, ok = parsed.Parsed["originName"]; !ok {
		return nil, fmt.Errorf("the segment 'originName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOriginID checks that 'input' can be parsed as a Origin ID
func ValidateOriginID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOriginID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Origin ID
func (id OriginId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/endpoints/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.OriginName)
}

// Segments returns a slice of Resource ID Segments which comprise this Origin ID
func (id OriginId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticOrigins", "origins", "origins"),
		resourceids.UserSpecifiedSegment("originName", "originValue"),
	}
}

// String returns a human-readable description of this Origin ID
func (id OriginId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
		fmt.Sprintf("Origin Name: %q", id.OriginName),
	}
	return fmt.Sprintf("Origin (%s)", strings.Join(components, "\n"))
}
