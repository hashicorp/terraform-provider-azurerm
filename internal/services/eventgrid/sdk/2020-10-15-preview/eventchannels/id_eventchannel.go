package eventchannels

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EventChannelId{}

// EventChannelId is a struct representing the Resource ID for a Event Channel
type EventChannelId struct {
	SubscriptionId       string
	ResourceGroupName    string
	PartnerNamespaceName string
	EventChannelName     string
}

// NewEventChannelID returns a new EventChannelId struct
func NewEventChannelID(subscriptionId string, resourceGroupName string, partnerNamespaceName string, eventChannelName string) EventChannelId {
	return EventChannelId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		PartnerNamespaceName: partnerNamespaceName,
		EventChannelName:     eventChannelName,
	}
}

// ParseEventChannelID parses 'input' into a EventChannelId
func ParseEventChannelID(input string) (*EventChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(EventChannelId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EventChannelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PartnerNamespaceName, ok = parsed.Parsed["partnerNamespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'partnerNamespaceName' was not found in the resource id %q", input)
	}

	if id.EventChannelName, ok = parsed.Parsed["eventChannelName"]; !ok {
		return nil, fmt.Errorf("the segment 'eventChannelName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseEventChannelIDInsensitively parses 'input' case-insensitively into a EventChannelId
// note: this method should only be used for API response data and not user input
func ParseEventChannelIDInsensitively(input string) (*EventChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(EventChannelId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EventChannelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PartnerNamespaceName, ok = parsed.Parsed["partnerNamespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'partnerNamespaceName' was not found in the resource id %q", input)
	}

	if id.EventChannelName, ok = parsed.Parsed["eventChannelName"]; !ok {
		return nil, fmt.Errorf("the segment 'eventChannelName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateEventChannelID checks that 'input' can be parsed as a Event Channel ID
func ValidateEventChannelID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEventChannelID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Event Channel ID
func (id EventChannelId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerNamespaces/%s/eventChannels/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerNamespaceName, id.EventChannelName)
}

// Segments returns a slice of Resource ID Segments which comprise this Event Channel ID
func (id EventChannelId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerNamespaces", "partnerNamespaces", "partnerNamespaces"),
		resourceids.UserSpecifiedSegment("partnerNamespaceName", "partnerNamespaceValue"),
		resourceids.StaticSegment("staticEventChannels", "eventChannels", "eventChannels"),
		resourceids.UserSpecifiedSegment("eventChannelName", "eventChannelValue"),
	}
}

// String returns a human-readable description of this Event Channel ID
func (id EventChannelId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Namespace Name: %q", id.PartnerNamespaceName),
		fmt.Sprintf("Event Channel Name: %q", id.EventChannelName),
	}
	return fmt.Sprintf("Event Channel (%s)", strings.Join(components, "\n"))
}
