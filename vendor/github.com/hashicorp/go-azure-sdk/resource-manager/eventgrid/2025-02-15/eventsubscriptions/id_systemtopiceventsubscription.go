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
	recaser.RegisterResourceId(&SystemTopicEventSubscriptionId{})
}

var _ resourceids.ResourceId = &SystemTopicEventSubscriptionId{}

// SystemTopicEventSubscriptionId is a struct representing the Resource ID for a System Topic Event Subscription
type SystemTopicEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	SystemTopicName       string
	EventSubscriptionName string
}

// NewSystemTopicEventSubscriptionID returns a new SystemTopicEventSubscriptionId struct
func NewSystemTopicEventSubscriptionID(subscriptionId string, resourceGroupName string, systemTopicName string, eventSubscriptionName string) SystemTopicEventSubscriptionId {
	return SystemTopicEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		SystemTopicName:       systemTopicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseSystemTopicEventSubscriptionID parses 'input' into a SystemTopicEventSubscriptionId
func ParseSystemTopicEventSubscriptionID(input string) (*SystemTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemTopicEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSystemTopicEventSubscriptionIDInsensitively parses 'input' case-insensitively into a SystemTopicEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseSystemTopicEventSubscriptionIDInsensitively(input string) (*SystemTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemTopicEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SystemTopicEventSubscriptionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SystemTopicName, ok = input.Parsed["systemTopicName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "systemTopicName", input)
	}

	if id.EventSubscriptionName, ok = input.Parsed["eventSubscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", input)
	}

	return nil
}

// ValidateSystemTopicEventSubscriptionID checks that 'input' can be parsed as a System Topic Event Subscription ID
func ValidateSystemTopicEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSystemTopicEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted System Topic Event Subscription ID
func (id SystemTopicEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/systemTopics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SystemTopicName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this System Topic Event Subscription ID
func (id SystemTopicEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticSystemTopics", "systemTopics", "systemTopics"),
		resourceids.UserSpecifiedSegment("systemTopicName", "systemTopicName"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionName"),
	}
}

// String returns a human-readable description of this System Topic Event Subscription ID
func (id SystemTopicEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("System Topic Name: %q", id.SystemTopicName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("System Topic Event Subscription (%s)", strings.Join(components, "\n"))
}
