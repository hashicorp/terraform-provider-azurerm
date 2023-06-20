package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = Subscriptions2Id{}

// Subscriptions2Id is a struct representing the Resource ID for a Subscriptions 2
type Subscriptions2Id struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	TopicName         string
	SubscriptionName  string
}

// NewSubscriptions2ID returns a new Subscriptions2Id struct
func NewSubscriptions2ID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string, subscriptionName string) Subscriptions2Id {
	return Subscriptions2Id{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		TopicName:         topicName,
		SubscriptionName:  subscriptionName,
	}
}

// ParseSubscriptions2ID parses 'input' into a Subscriptions2Id
func ParseSubscriptions2ID(input string) (*Subscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(Subscriptions2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Subscriptions2Id{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", *parsed)
	}

	return &id, nil
}

// ParseSubscriptions2IDInsensitively parses 'input' case-insensitively into a Subscriptions2Id
// note: this method should only be used for API response data and not user input
func ParseSubscriptions2IDInsensitively(input string) (*Subscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(Subscriptions2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Subscriptions2Id{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", *parsed)
	}

	return &id, nil
}

// ValidateSubscriptions2ID checks that 'input' can be parsed as a Subscriptions 2 ID
func ValidateSubscriptions2ID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubscriptions2ID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subscriptions 2 ID
func (id Subscriptions2Id) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.SubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subscriptions 2 ID
func (id Subscriptions2Id) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticSubscriptions2", "subscriptions", "subscriptions"),
		resourceids.UserSpecifiedSegment("subscriptionName", "subscriptionValue"),
	}
}

// String returns a human-readable description of this Subscriptions 2 ID
func (id Subscriptions2Id) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Subscription Name: %q", id.SubscriptionName),
	}
	return fmt.Sprintf("Subscriptions 2 (%s)", strings.Join(components, "\n"))
}
