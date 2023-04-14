package actiongroupsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NotificationStatusId{}

// NotificationStatusId is a struct representing the Resource ID for a Notification Status
type NotificationStatusId struct {
	SubscriptionId string
	NotificationId string
}

// NewNotificationStatusID returns a new NotificationStatusId struct
func NewNotificationStatusID(subscriptionId string, notificationId string) NotificationStatusId {
	return NotificationStatusId{
		SubscriptionId: subscriptionId,
		NotificationId: notificationId,
	}
}

// ParseNotificationStatusID parses 'input' into a NotificationStatusId
func ParseNotificationStatusID(input string) (*NotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationStatusId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseNotificationStatusIDInsensitively parses 'input' case-insensitively into a NotificationStatusId
// note: this method should only be used for API response data and not user input
func ParseNotificationStatusIDInsensitively(input string) (*NotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationStatusId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateNotificationStatusID checks that 'input' can be parsed as a Notification Status ID
func ValidateNotificationStatusID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationStatusID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification Status ID
func (id NotificationStatusId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Insights/notificationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.NotificationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Notification Status ID
func (id NotificationStatusId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticNotificationStatus", "notificationStatus", "notificationStatus"),
		resourceids.UserSpecifiedSegment("notificationId", "notificationIdValue"),
	}
}

// String returns a human-readable description of this Notification Status ID
func (id NotificationStatusId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Notification: %q", id.NotificationId),
	}
	return fmt.Sprintf("Notification Status (%s)", strings.Join(components, "\n"))
}
