package consumergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConsumerGroupId{}

// ConsumerGroupId is a struct representing the Resource ID for a Consumer Group
type ConsumerGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	EventHubName      string
	ConsumerGroupName string
}

// NewConsumerGroupID returns a new ConsumerGroupId struct
func NewConsumerGroupID(subscriptionId string, resourceGroupName string, namespaceName string, eventHubName string, consumerGroupName string) ConsumerGroupId {
	return ConsumerGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		EventHubName:      eventHubName,
		ConsumerGroupName: consumerGroupName,
	}
}

// ParseConsumerGroupID parses 'input' into a ConsumerGroupId
func ParseConsumerGroupID(input string) (*ConsumerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConsumerGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConsumerGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.EventHubName, ok = parsed.Parsed["eventHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'eventHubName' was not found in the resource id %q", input)
	}

	if id.ConsumerGroupName, ok = parsed.Parsed["consumerGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'consumerGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConsumerGroupIDInsensitively parses 'input' case-insensitively into a ConsumerGroupId
// note: this method should only be used for API response data and not user input
func ParseConsumerGroupIDInsensitively(input string) (*ConsumerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConsumerGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConsumerGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.EventHubName, ok = parsed.Parsed["eventHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'eventHubName' was not found in the resource id %q", input)
	}

	if id.ConsumerGroupName, ok = parsed.Parsed["consumerGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'consumerGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConsumerGroupID checks that 'input' can be parsed as a Consumer Group ID
func ValidateConsumerGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConsumerGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Consumer Group ID
func (id ConsumerGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/consumerGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventHubName, id.ConsumerGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Consumer Group ID
func (id ConsumerGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticEventhubs", "eventhubs", "eventhubs"),
		resourceids.UserSpecifiedSegment("eventHubName", "eventHubValue"),
		resourceids.StaticSegment("staticConsumerGroups", "consumerGroups", "consumerGroups"),
		resourceids.UserSpecifiedSegment("consumerGroupName", "consumerGroupValue"),
	}
}

// String returns a human-readable description of this Consumer Group ID
func (id ConsumerGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Event Hub Name: %q", id.EventHubName),
		fmt.Sprintf("Consumer Group Name: %q", id.ConsumerGroupName),
	}
	return fmt.Sprintf("Consumer Group (%s)", strings.Join(components, "\n"))
}
