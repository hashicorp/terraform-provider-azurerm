package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PartnerTopicEventSubscriptionId{}

// PartnerTopicEventSubscriptionId is a struct representing the Resource ID for a Partner Topic Event Subscription
type PartnerTopicEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	PartnerTopicName      string
	EventSubscriptionName string
}

// NewPartnerTopicEventSubscriptionID returns a new PartnerTopicEventSubscriptionId struct
func NewPartnerTopicEventSubscriptionID(subscriptionId string, resourceGroupName string, partnerTopicName string, eventSubscriptionName string) PartnerTopicEventSubscriptionId {
	return PartnerTopicEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		PartnerTopicName:      partnerTopicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParsePartnerTopicEventSubscriptionID parses 'input' into a PartnerTopicEventSubscriptionId
func ParsePartnerTopicEventSubscriptionID(input string) (*PartnerTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerTopicEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerTopicName, ok = parsed.Parsed["partnerTopicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerTopicName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ParsePartnerTopicEventSubscriptionIDInsensitively parses 'input' case-insensitively into a PartnerTopicEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParsePartnerTopicEventSubscriptionIDInsensitively(input string) (*PartnerTopicEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerTopicEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerTopicEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerTopicName, ok = parsed.Parsed["partnerTopicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerTopicName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ValidatePartnerTopicEventSubscriptionID checks that 'input' can be parsed as a Partner Topic Event Subscription ID
func ValidatePartnerTopicEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartnerTopicEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partner Topic Event Subscription ID
func (id PartnerTopicEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerTopics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerTopicName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Partner Topic Event Subscription ID
func (id PartnerTopicEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerTopics", "partnerTopics", "partnerTopics"),
		resourceids.UserSpecifiedSegment("partnerTopicName", "partnerTopicValue"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionValue"),
	}
}

// String returns a human-readable description of this Partner Topic Event Subscription ID
func (id PartnerTopicEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Topic Name: %q", id.PartnerTopicName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Partner Topic Event Subscription (%s)", strings.Join(components, "\n"))
}
