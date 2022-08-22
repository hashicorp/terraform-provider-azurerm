package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TopicId{}

// TopicId is a struct representing the Resource ID for a Topic
type TopicId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	TopicName         string
}

// NewTopicID returns a new TopicId struct
func NewTopicID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string) TopicId {
	return TopicId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		TopicName:         topicName,
	}
}

// ParseTopicID parses 'input' into a TopicId
func ParseTopicID(input string) (*TopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, fmt.Errorf("the segment 'topicName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTopicIDInsensitively parses 'input' case-insensitively into a TopicId
// note: this method should only be used for API response data and not user input
func ParseTopicIDInsensitively(input string) (*TopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, fmt.Errorf("the segment 'topicName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTopicID checks that 'input' can be parsed as a Topic ID
func ValidateTopicID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTopicID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Topic ID
func (id TopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName)
}

// Segments returns a slice of Resource ID Segments which comprise this Topic ID
func (id TopicId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicValue"),
	}
}

// String returns a human-readable description of this Topic ID
func (id TopicId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
	}
	return fmt.Sprintf("Topic (%s)", strings.Join(components, "\n"))
}
