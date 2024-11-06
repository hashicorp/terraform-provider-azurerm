package hubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NotificationHubId{})
}

var _ resourceids.ResourceId = &NotificationHubId{}

// NotificationHubId is a struct representing the Resource ID for a Notification Hub
type NotificationHubId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NamespaceName       string
	NotificationHubName string
}

// NewNotificationHubID returns a new NotificationHubId struct
func NewNotificationHubID(subscriptionId string, resourceGroupName string, namespaceName string, notificationHubName string) NotificationHubId {
	return NotificationHubId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NamespaceName:       namespaceName,
		NotificationHubName: notificationHubName,
	}
}

// ParseNotificationHubID parses 'input' into a NotificationHubId
func ParseNotificationHubID(input string) (*NotificationHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationHubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNotificationHubIDInsensitively parses 'input' case-insensitively into a NotificationHubId
// note: this method should only be used for API response data and not user input
func ParseNotificationHubIDInsensitively(input string) (*NotificationHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationHubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NotificationHubId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.NotificationHubName, ok = input.Parsed["notificationHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "notificationHubName", input)
	}

	return nil
}

// ValidateNotificationHubID checks that 'input' can be parsed as a Notification Hub ID
func ValidateNotificationHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification Hub ID
func (id NotificationHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.NotificationHubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Notification Hub ID
func (id NotificationHubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNotificationHubs", "Microsoft.NotificationHubs", "Microsoft.NotificationHubs"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticNotificationHubs", "notificationHubs", "notificationHubs"),
		resourceids.UserSpecifiedSegment("notificationHubName", "notificationHubName"),
	}
}

// String returns a human-readable description of this Notification Hub ID
func (id NotificationHubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Notification Hub Name: %q", id.NotificationHubName),
	}
	return fmt.Sprintf("Notification Hub (%s)", strings.Join(components, "\n"))
}
