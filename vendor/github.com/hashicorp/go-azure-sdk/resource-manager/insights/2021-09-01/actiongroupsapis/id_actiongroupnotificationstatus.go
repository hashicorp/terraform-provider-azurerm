package actiongroupsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ActionGroupNotificationStatusId{}

// ActionGroupNotificationStatusId is a struct representing the Resource ID for a Action Group Notification Status
type ActionGroupNotificationStatusId struct {
	SubscriptionId    string
	ResourceGroupName string
	ActionGroupName   string
	NotificationId    string
}

// NewActionGroupNotificationStatusID returns a new ActionGroupNotificationStatusId struct
func NewActionGroupNotificationStatusID(subscriptionId string, resourceGroupName string, actionGroupName string, notificationId string) ActionGroupNotificationStatusId {
	return ActionGroupNotificationStatusId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ActionGroupName:   actionGroupName,
		NotificationId:    notificationId,
	}
}

// ParseActionGroupNotificationStatusID parses 'input' into a ActionGroupNotificationStatusId
func ParseActionGroupNotificationStatusID(input string) (*ActionGroupNotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActionGroupNotificationStatusId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActionGroupNotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ActionGroupName, ok = parsed.Parsed["actionGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'actionGroupName' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseActionGroupNotificationStatusIDInsensitively parses 'input' case-insensitively into a ActionGroupNotificationStatusId
// note: this method should only be used for API response data and not user input
func ParseActionGroupNotificationStatusIDInsensitively(input string) (*ActionGroupNotificationStatusId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActionGroupNotificationStatusId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActionGroupNotificationStatusId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ActionGroupName, ok = parsed.Parsed["actionGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'actionGroupName' was not found in the resource id %q", input)
	}

	if id.NotificationId, ok = parsed.Parsed["notificationId"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateActionGroupNotificationStatusID checks that 'input' can be parsed as a Action Group Notification Status ID
func ValidateActionGroupNotificationStatusID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActionGroupNotificationStatusID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Action Group Notification Status ID
func (id ActionGroupNotificationStatusId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/actionGroups/%s/notificationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ActionGroupName, id.NotificationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Action Group Notification Status ID
func (id ActionGroupNotificationStatusId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticActionGroups", "actionGroups", "actionGroups"),
		resourceids.UserSpecifiedSegment("actionGroupName", "actionGroupValue"),
		resourceids.StaticSegment("staticNotificationStatus", "notificationStatus", "notificationStatus"),
		resourceids.UserSpecifiedSegment("notificationId", "notificationIdValue"),
	}
}

// String returns a human-readable description of this Action Group Notification Status ID
func (id ActionGroupNotificationStatusId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Action Group Name: %q", id.ActionGroupName),
		fmt.Sprintf("Notification: %q", id.NotificationId),
	}
	return fmt.Sprintf("Action Group Notification Status (%s)", strings.Join(components, "\n"))
}
