package consumergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConsumerGroupId{})
}

var _ resourceids.ResourceId = &ConsumerGroupId{}

// ConsumerGroupId is a struct representing the Resource ID for a Consumer Group
type ConsumerGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	EventhubName      string
	ConsumerGroupName string
}

// NewConsumerGroupID returns a new ConsumerGroupId struct
func NewConsumerGroupID(subscriptionId string, resourceGroupName string, namespaceName string, eventhubName string, consumerGroupName string) ConsumerGroupId {
	return ConsumerGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		EventhubName:      eventhubName,
		ConsumerGroupName: consumerGroupName,
	}
}

// ParseConsumerGroupID parses 'input' into a ConsumerGroupId
func ParseConsumerGroupID(input string) (*ConsumerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConsumerGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConsumerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConsumerGroupIDInsensitively parses 'input' case-insensitively into a ConsumerGroupId
// note: this method should only be used for API response data and not user input
func ParseConsumerGroupIDInsensitively(input string) (*ConsumerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConsumerGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConsumerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConsumerGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.EventhubName, ok = input.Parsed["eventhubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventhubName", input)
	}

	if id.ConsumerGroupName, ok = input.Parsed["consumerGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "consumerGroupName", input)
	}

	return nil
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName, id.ConsumerGroupName)
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
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticEventhubs", "eventhubs", "eventhubs"),
		resourceids.UserSpecifiedSegment("eventhubName", "eventhubName"),
		resourceids.StaticSegment("staticConsumerGroups", "consumerGroups", "consumerGroups"),
		resourceids.UserSpecifiedSegment("consumerGroupName", "consumerGroupName"),
	}
}

// String returns a human-readable description of this Consumer Group ID
func (id ConsumerGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Eventhub Name: %q", id.EventhubName),
		fmt.Sprintf("Consumer Group Name: %q", id.ConsumerGroupName),
	}
	return fmt.Sprintf("Consumer Group (%s)", strings.Join(components, "\n"))
}
