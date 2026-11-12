package eventhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EventhubId{})
}

var _ resourceids.ResourceId = &EventhubId{}

// EventhubId is a struct representing the Resource ID for a Eventhub
type EventhubId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	EventhubName      string
}

// NewEventhubID returns a new EventhubId struct
func NewEventhubID(subscriptionId string, resourceGroupName string, namespaceName string, eventhubName string) EventhubId {
	return EventhubId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		EventhubName:      eventhubName,
	}
}

// ParseEventhubID parses 'input' into a EventhubId
func ParseEventhubID(input string) (*EventhubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventhubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventhubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEventhubIDInsensitively parses 'input' case-insensitively into a EventhubId
// note: this method should only be used for API response data and not user input
func ParseEventhubIDInsensitively(input string) (*EventhubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventhubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventhubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EventhubId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateEventhubID checks that 'input' can be parsed as a Eventhub ID
func ValidateEventhubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEventhubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Eventhub ID
func (id EventhubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Eventhub ID
func (id EventhubId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Eventhub ID
func (id EventhubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Eventhub Name: %q", id.EventhubName),
	}
	return fmt.Sprintf("Eventhub (%s)", strings.Join(components, "\n"))
}
