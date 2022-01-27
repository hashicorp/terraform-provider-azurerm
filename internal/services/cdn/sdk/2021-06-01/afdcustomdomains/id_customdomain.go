package afdcustomdomains

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CustomDomainId{}

// CustomDomainId is a struct representing the Resource ID for a Custom Domain
type CustomDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	CustomDomainName  string
}

// NewCustomDomainID returns a new CustomDomainId struct
func NewCustomDomainID(subscriptionId string, resourceGroupName string, profileName string, customDomainName string) CustomDomainId {
	return CustomDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		CustomDomainName:  customDomainName,
	}
}

// ParseCustomDomainID parses 'input' into a CustomDomainId
func ParseCustomDomainID(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomDomainId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.CustomDomainName, ok = parsed.Parsed["customDomainName"]; !ok {
		return nil, fmt.Errorf("the segment 'customDomainName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCustomDomainIDInsensitively parses 'input' case-insensitively into a CustomDomainId
// note: this method should only be used for API response data and not user input
func ParseCustomDomainIDInsensitively(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomDomainId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.CustomDomainName, ok = parsed.Parsed["customDomainName"]; !ok {
		return nil, fmt.Errorf("the segment 'customDomainName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCustomDomainID checks that 'input' can be parsed as a Custom Domain ID
func ValidateCustomDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom Domain ID
func (id CustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.CustomDomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom Domain ID
func (id CustomDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticCustomDomains", "customDomains", "customDomains"),
		resourceids.UserSpecifiedSegment("customDomainName", "customDomainValue"),
	}
}

// String returns a human-readable description of this Custom Domain ID
func (id CustomDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Custom Domain Name: %q", id.CustomDomainName),
	}
	return fmt.Sprintf("Custom Domain (%s)", strings.Join(components, "\n"))
}
