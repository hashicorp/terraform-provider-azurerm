package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NamespaceTopicEventSubscriptionId{})
}

var _ resourceids.ResourceId = &NamespaceTopicEventSubscriptionId{}

// NamespaceTopicEventSubscriptionId is a struct representing the Resource ID for a Namespace Topic Event Subscription
type NamespaceTopicEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	TopicName             string
	EventSubscriptionName string
}

// NewNamespaceTopicEventSubscriptionID returns a new NamespaceTopicEventSubscriptionId struct
func NewNamespaceTopicEventSubscriptionID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string, eventSubscriptionName string) NamespaceTopicEventSubscriptionId {
	return NamespaceTopicEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		TopicName:             topicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseNamespaceTopicEventSubscriptionID parses 'input' into a NamespaceTopicEventSubscriptionId
func ParseNamespaceTopicEventSubscriptionID(input string) (*NamespaceTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceTopicEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNamespaceTopicEventSubscriptionIDInsensitively parses 'input' case-insensitively into a NamespaceTopicEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseNamespaceTopicEventSubscriptionIDInsensitively(input string) (*NamespaceTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceTopicEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NamespaceTopicEventSubscriptionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TopicName, ok = input.Parsed["topicName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "topicName", input)
	}

	if id.EventSubscriptionName, ok = input.Parsed["eventSubscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", input)
	}

	return nil
}

// ValidateNamespaceTopicEventSubscriptionID checks that 'input' can be parsed as a Namespace Topic Event Subscription ID
func ValidateNamespaceTopicEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNamespaceTopicEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Namespace Topic Event Subscription ID
func (id NamespaceTopicEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/topics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Namespace Topic Event Subscription ID
func (id NamespaceTopicEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicName"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionName"),
	}
}

// String returns a human-readable description of this Namespace Topic Event Subscription ID
func (id NamespaceTopicEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Namespace Topic Event Subscription (%s)", strings.Join(components, "\n"))
}
