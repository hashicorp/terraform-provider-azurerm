package queuesauthorizationrule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = QueueId{}

// QueueId is a struct representing the Resource ID for a Queue
type QueueId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	QueueName         string
}

// NewQueueID returns a new QueueId struct
func NewQueueID(subscriptionId string, resourceGroupName string, namespaceName string, queueName string) QueueId {
	return QueueId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		QueueName:         queueName,
	}
}

// ParseQueueID parses 'input' into a QueueId
func ParseQueueID(input string) (*QueueId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueueId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.QueueName, ok = parsed.Parsed["queueName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queueName", *parsed)
	}

	return &id, nil
}

// ParseQueueIDInsensitively parses 'input' case-insensitively into a QueueId
// note: this method should only be used for API response data and not user input
func ParseQueueIDInsensitively(input string) (*QueueId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueueId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.QueueName, ok = parsed.Parsed["queueName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queueName", *parsed)
	}

	return &id, nil
}

// ValidateQueueID checks that 'input' can be parsed as a Queue ID
func ValidateQueueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQueueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Queue ID
func (id QueueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/queues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.QueueName)
}

// Segments returns a slice of Resource ID Segments which comprise this Queue ID
func (id QueueId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticQueues", "queues", "queues"),
		resourceids.UserSpecifiedSegment("queueName", "queueValue"),
	}
}

// String returns a human-readable description of this Queue ID
func (id QueueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Queue Name: %q", id.QueueName),
	}
	return fmt.Sprintf("Queue (%s)", strings.Join(components, "\n"))
}
