package streamingpoliciesandstreaminglocators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = StreamingLocatorId{}

// StreamingLocatorId is a struct representing the Resource ID for a Streaming Locator
type StreamingLocatorId struct {
	SubscriptionId       string
	ResourceGroupName    string
	AccountName          string
	StreamingLocatorName string
}

// NewStreamingLocatorID returns a new StreamingLocatorId struct
func NewStreamingLocatorID(subscriptionId string, resourceGroupName string, accountName string, streamingLocatorName string) StreamingLocatorId {
	return StreamingLocatorId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		AccountName:          accountName,
		StreamingLocatorName: streamingLocatorName,
	}
}

// ParseStreamingLocatorID parses 'input' into a StreamingLocatorId
func ParseStreamingLocatorID(input string) (*StreamingLocatorId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingLocatorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingLocatorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.StreamingLocatorName, ok = parsed.Parsed["streamingLocatorName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingLocatorName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseStreamingLocatorIDInsensitively parses 'input' case-insensitively into a StreamingLocatorId
// note: this method should only be used for API response data and not user input
func ParseStreamingLocatorIDInsensitively(input string) (*StreamingLocatorId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingLocatorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingLocatorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.StreamingLocatorName, ok = parsed.Parsed["streamingLocatorName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingLocatorName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateStreamingLocatorID checks that 'input' can be parsed as a Streaming Locator ID
func ValidateStreamingLocatorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStreamingLocatorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Streaming Locator ID
func (id StreamingLocatorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/streamingLocators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.StreamingLocatorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Streaming Locator ID
func (id StreamingLocatorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticStreamingLocators", "streamingLocators", "streamingLocators"),
		resourceids.UserSpecifiedSegment("streamingLocatorName", "streamingLocatorValue"),
	}
}

// String returns a human-readable description of this Streaming Locator ID
func (id StreamingLocatorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Streaming Locator Name: %q", id.StreamingLocatorName),
	}
	return fmt.Sprintf("Streaming Locator (%s)", strings.Join(components, "\n"))
}
