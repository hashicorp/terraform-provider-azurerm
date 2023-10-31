package notificationrecipientemail

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NotificationId{}

// NotificationId is a struct representing the Resource ID for a Notification
type NotificationId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	NotificationName  NotificationName
}

// NewNotificationID returns a new NotificationId struct
func NewNotificationID(subscriptionId string, resourceGroupName string, serviceName string, notificationName NotificationName) NotificationId {
	return NotificationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		NotificationName:  notificationName,
	}
}

// ParseNotificationID parses 'input' into a NotificationId
func ParseNotificationID(input string) (*NotificationId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if v, ok := parsed.Parsed["notificationName"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "notificationName", *parsed)
		}

		notificationName, err := parseNotificationName(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.NotificationName = *notificationName
	}

	return &id, nil
}

// ParseNotificationIDInsensitively parses 'input' case-insensitively into a NotificationId
// note: this method should only be used for API response data and not user input
func ParseNotificationIDInsensitively(input string) (*NotificationId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if v, ok := parsed.Parsed["notificationName"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "notificationName", *parsed)
		}

		notificationName, err := parseNotificationName(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.NotificationName = *notificationName
	}

	return &id, nil
}

// ValidateNotificationID checks that 'input' can be parsed as a Notification ID
func ValidateNotificationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification ID
func (id NotificationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/notifications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, string(id.NotificationName))
}

// Segments returns a slice of Resource ID Segments which comprise this Notification ID
func (id NotificationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticNotifications", "notifications", "notifications"),
		resourceids.ConstantSegment("notificationName", PossibleValuesForNotificationName(), "AccountClosedPublisher"),
	}
}

// String returns a human-readable description of this Notification ID
func (id NotificationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Notification Name: %q", string(id.NotificationName)),
	}
	return fmt.Sprintf("Notification (%s)", strings.Join(components, "\n"))
}
