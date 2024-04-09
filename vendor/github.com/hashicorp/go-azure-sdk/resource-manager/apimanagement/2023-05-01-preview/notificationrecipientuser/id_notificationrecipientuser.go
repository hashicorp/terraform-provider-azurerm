package notificationrecipientuser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &NotificationRecipientUserId{}

// NotificationRecipientUserId is a struct representing the Resource ID for a Notification Recipient User
type NotificationRecipientUserId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	NotificationName  NotificationName
	UserId            string
}

// NewNotificationRecipientUserID returns a new NotificationRecipientUserId struct
func NewNotificationRecipientUserID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, notificationName NotificationName, userId string) NotificationRecipientUserId {
	return NotificationRecipientUserId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		NotificationName:  notificationName,
		UserId:            userId,
	}
}

// ParseNotificationRecipientUserID parses 'input' into a NotificationRecipientUserId
func ParseNotificationRecipientUserID(input string) (*NotificationRecipientUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationRecipientUserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationRecipientUserId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNotificationRecipientUserIDInsensitively parses 'input' case-insensitively into a NotificationRecipientUserId
// note: this method should only be used for API response data and not user input
func ParseNotificationRecipientUserIDInsensitively(input string) (*NotificationRecipientUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationRecipientUserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationRecipientUserId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NotificationRecipientUserId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.WorkspaceId, ok = input.Parsed["workspaceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceId", input)
	}

	if v, ok := input.Parsed["notificationName"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "notificationName", input)
		}

		notificationName, err := parseNotificationName(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.NotificationName = *notificationName
	}

	if id.UserId, ok = input.Parsed["userId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "userId", input)
	}

	return nil
}

// ValidateNotificationRecipientUserID checks that 'input' can be parsed as a Notification Recipient User ID
func ValidateNotificationRecipientUserID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationRecipientUserID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification Recipient User ID
func (id NotificationRecipientUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/notifications/%s/recipientUsers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, string(id.NotificationName), id.UserId)
}

// Segments returns a slice of Resource ID Segments which comprise this Notification Recipient User ID
func (id NotificationRecipientUserId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceIdValue"),
		resourceids.StaticSegment("staticNotifications", "notifications", "notifications"),
		resourceids.ConstantSegment("notificationName", PossibleValuesForNotificationName(), "AccountClosedPublisher"),
		resourceids.StaticSegment("staticRecipientUsers", "recipientUsers", "recipientUsers"),
		resourceids.UserSpecifiedSegment("userId", "userIdValue"),
	}
}

// String returns a human-readable description of this Notification Recipient User ID
func (id NotificationRecipientUserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Notification Name: %q", string(id.NotificationName)),
		fmt.Sprintf("User: %q", id.UserId),
	}
	return fmt.Sprintf("Notification Recipient User (%s)", strings.Join(components, "\n"))
}
