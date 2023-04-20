package actiongroupsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderNotificationStatusId{}

// ProviderNotificationStatusId is a struct representing the Resource ID for a Provider Notification Status
type ProviderNotificationStatusId struct {
	SubscriptionId    string
	ResourceGroupName string
	NotificationId    string
}

// NewProviderNotificationStatusID returns a new ProviderNotificationStatusId struct
func NewProviderNotificationStatusID(subscriptionId string, resourceGroupName string, notificationId string) ProviderNotificationStatusId {
	return ProviderNotificationStatusId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NotificationId:    notificationId,
	}
}

// ParseProviderNotificationStatusID parses 'input' into a ProviderNotificationStatusId
func ParseProviderNotificationStatusID(input string) (*ProviderNotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderNotificationStatusId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderNotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProviderNotificationStatusIDInsensitively parses 'input' case-insensitively into a ProviderNotificationStatusId
// note: this method should only be used for API response data and not user input
func ParseProviderNotificationStatusIDInsensitively(input string) (*ProviderNotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderNotificationStatusId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderNotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProviderNotificationStatusID checks that 'input' can be parsed as a Provider Notification Status ID
func ValidateProviderNotificationStatusID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderNotificationStatusID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Notification Status ID
func (id ProviderNotificationStatusId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/notificationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NotificationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Notification Status ID
func (id ProviderNotificationStatusId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticNotificationStatus", "notificationStatus", "notificationStatus"),
		resourceids.UserSpecifiedSegment("notificationId", "notificationIdValue"),
	}
}

// String returns a human-readable description of this Provider Notification Status ID
func (id ProviderNotificationStatusId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Notification: %q", id.NotificationId),
	}
	return fmt.Sprintf("Provider Notification Status (%s)", strings.Join(components, "\n"))
}
