package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TopicEventSubscriptionId{}

// TopicEventSubscriptionId is a struct representing the Resource ID for a Topic Event Subscription
type TopicEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DomainName            string
	TopicName             string
	EventSubscriptionName string
}

// NewTopicEventSubscriptionID returns a new TopicEventSubscriptionId struct
func NewTopicEventSubscriptionID(subscriptionId string, resourceGroupName string, domainName string, topicName string, eventSubscriptionName string) TopicEventSubscriptionId {
	return TopicEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DomainName:            domainName,
		TopicName:             topicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseTopicEventSubscriptionID parses 'input' into a TopicEventSubscriptionId
func ParseTopicEventSubscriptionID(input string) (*TopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ParseTopicEventSubscriptionIDInsensitively parses 'input' case-insensitively into a TopicEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseTopicEventSubscriptionIDInsensitively(input string) (*TopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ValidateTopicEventSubscriptionID checks that 'input' can be parsed as a Topic Event Subscription ID
func ValidateTopicEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTopicEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Topic Event Subscription ID
func (id TopicEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/domains/%s/topics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DomainName, id.TopicName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Topic Event Subscription ID
func (id TopicEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainValue"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicValue"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionValue"),
	}
}

// String returns a human-readable description of this Topic Event Subscription ID
func (id TopicEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Topic Event Subscription (%s)", strings.Join(components, "\n"))
}
