package accounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProviderMediaServiceId{}

// ProviderMediaServiceId is a struct representing the Resource ID for a Provider Media Service
type ProviderMediaServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
}

// NewProviderMediaServiceID returns a new ProviderMediaServiceId struct
func NewProviderMediaServiceID(subscriptionId string, resourceGroupName string, mediaServiceName string) ProviderMediaServiceId {
	return ProviderMediaServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
	}
}

// ParseProviderMediaServiceID parses 'input' into a ProviderMediaServiceId
func ParseProviderMediaServiceID(input string) (*ProviderMediaServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderMediaServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderMediaServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'mediaServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProviderMediaServiceIDInsensitively parses 'input' case-insensitively into a ProviderMediaServiceId
// note: this method should only be used for API response data and not user input
func ParseProviderMediaServiceIDInsensitively(input string) (*ProviderMediaServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderMediaServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderMediaServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'mediaServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProviderMediaServiceID checks that 'input' can be parsed as a Provider Media Service ID
func ValidateProviderMediaServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderMediaServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Media Service ID
func (id ProviderMediaServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Media Service ID
func (id ProviderMediaServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
	}
}

// String returns a human-readable description of this Provider Media Service ID
func (id ProviderMediaServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
	}
	return fmt.Sprintf("Provider Media Service (%s)", strings.Join(components, "\n"))
}
