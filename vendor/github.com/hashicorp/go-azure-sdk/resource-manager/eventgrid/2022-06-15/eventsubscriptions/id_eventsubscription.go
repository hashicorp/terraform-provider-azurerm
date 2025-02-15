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
	recaser.RegisterResourceId(&EventSubscriptionId{})
}

var _ resourceids.ResourceId = &EventSubscriptionId{}

// EventSubscriptionId is a struct representing the Resource ID for a Event Subscription
type EventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	TopicName             string
	EventSubscriptionName string
}

// NewEventSubscriptionID returns a new EventSubscriptionId struct
func NewEventSubscriptionID(subscriptionId string, resourceGroupName string, topicName string, eventSubscriptionName string) EventSubscriptionId {
	return EventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		TopicName:             topicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseEventSubscriptionID parses 'input' into a EventSubscriptionId
func ParseEventSubscriptionID(input string) (*EventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEventSubscriptionIDInsensitively parses 'input' case-insensitively into a EventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseEventSubscriptionIDInsensitively(input string) (*EventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EventSubscriptionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.TopicName, ok = input.Parsed["topicName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "topicName", input)
	}

	if id.EventSubscriptionName, ok = input.Parsed["eventSubscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", input)
	}

	return nil
}

// ValidateEventSubscriptionID checks that 'input' can be parsed as a Event Subscription ID
func ValidateEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Event Subscription ID
func (id EventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/topics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TopicName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Event Subscription ID
func (id EventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicName"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionName"),
	}
}

// String returns a human-readable description of this Event Subscription ID
func (id EventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Event Subscription (%s)", strings.Join(components, "\n"))
}
