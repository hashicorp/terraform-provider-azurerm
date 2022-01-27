package afdorigins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OriginGroupOriginId{}

// OriginGroupOriginId is a struct representing the Resource ID for a Origin Group Origin
type OriginGroupOriginId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	OriginGroupName   string
	OriginName        string
}

// NewOriginGroupOriginID returns a new OriginGroupOriginId struct
func NewOriginGroupOriginID(subscriptionId string, resourceGroupName string, profileName string, originGroupName string, originName string) OriginGroupOriginId {
	return OriginGroupOriginId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		OriginGroupName:   originGroupName,
		OriginName:        originName,
	}
}

// ParseOriginGroupOriginID parses 'input' into a OriginGroupOriginId
func ParseOriginGroupOriginID(input string) (*OriginGroupOriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(OriginGroupOriginId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OriginGroupOriginId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.OriginGroupName, ok = parsed.Parsed["originGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'originGroupName' was not found in the resource id %q", input)
	}

	if id.OriginName, ok = parsed.Parsed["originName"]; !ok {
		return nil, fmt.Errorf("the segment 'originName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOriginGroupOriginIDInsensitively parses 'input' case-insensitively into a OriginGroupOriginId
// note: this method should only be used for API response data and not user input
func ParseOriginGroupOriginIDInsensitively(input string) (*OriginGroupOriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(OriginGroupOriginId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OriginGroupOriginId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.OriginGroupName, ok = parsed.Parsed["originGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'originGroupName' was not found in the resource id %q", input)
	}

	if id.OriginName, ok = parsed.Parsed["originName"]; !ok {
		return nil, fmt.Errorf("the segment 'originName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOriginGroupOriginID checks that 'input' can be parsed as a Origin Group Origin ID
func ValidateOriginGroupOriginID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOriginGroupOriginID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Origin Group Origin ID
func (id OriginGroupOriginId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/originGroups/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.OriginGroupName, id.OriginName)
}

// Segments returns a slice of Resource ID Segments which comprise this Origin Group Origin ID
func (id OriginGroupOriginId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticOriginGroups", "originGroups", "originGroups"),
		resourceids.UserSpecifiedSegment("originGroupName", "originGroupValue"),
		resourceids.StaticSegment("staticOrigins", "origins", "origins"),
		resourceids.UserSpecifiedSegment("originName", "originValue"),
	}
}

// String returns a human-readable description of this Origin Group Origin ID
func (id OriginGroupOriginId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Origin Group Name: %q", id.OriginGroupName),
		fmt.Sprintf("Origin Name: %q", id.OriginName),
	}
	return fmt.Sprintf("Origin Group Origin (%s)", strings.Join(components, "\n"))
}
