package profiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TrafficManagerProfileId{}

// TrafficManagerProfileId is a struct representing the Resource ID for a Traffic Manager Profile
type TrafficManagerProfileId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
}

// NewTrafficManagerProfileID returns a new TrafficManagerProfileId struct
func NewTrafficManagerProfileID(subscriptionId string, resourceGroupName string, profileName string) TrafficManagerProfileId {
	return TrafficManagerProfileId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
	}
}

// ParseTrafficManagerProfileID parses 'input' into a TrafficManagerProfileId
func ParseTrafficManagerProfileID(input string) (*TrafficManagerProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrafficManagerProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrafficManagerProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTrafficManagerProfileIDInsensitively parses 'input' case-insensitively into a TrafficManagerProfileId
// note: this method should only be used for API response data and not user input
func ParseTrafficManagerProfileIDInsensitively(input string) (*TrafficManagerProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrafficManagerProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrafficManagerProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTrafficManagerProfileID checks that 'input' can be parsed as a Traffic Manager Profile ID
func ValidateTrafficManagerProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrafficManagerProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Traffic Manager Profile ID
func (id TrafficManagerProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Traffic Manager Profile ID
func (id TrafficManagerProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticTrafficManagerProfiles", "trafficManagerProfiles", "trafficManagerProfiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
	}
}

// String returns a human-readable description of this Traffic Manager Profile ID
func (id TrafficManagerProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
	}
	return fmt.Sprintf("Traffic Manager Profile (%s)", strings.Join(components, "\n"))
}
