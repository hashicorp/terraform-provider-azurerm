package notificationrecipientuser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RecipientUserId{})
}

var _ resourceids.ResourceId = &RecipientUserId{}

// RecipientUserId is a struct representing the Resource ID for a Recipient User
type RecipientUserId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	NotificationName  NotificationName
	UserId            string
}

// NewRecipientUserID returns a new RecipientUserId struct
func NewRecipientUserID(subscriptionId string, resourceGroupName string, serviceName string, notificationName NotificationName, userId string) RecipientUserId {
	return RecipientUserId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		NotificationName:  notificationName,
		UserId:            userId,
	}
}

// ParseRecipientUserID parses 'input' into a RecipientUserId
func ParseRecipientUserID(input string) (*RecipientUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RecipientUserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RecipientUserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRecipientUserIDInsensitively parses 'input' case-insensitively into a RecipientUserId
// note: this method should only be used for API response data and not user input
func ParseRecipientUserIDInsensitively(input string) (*RecipientUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RecipientUserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RecipientUserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RecipientUserId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateRecipientUserID checks that 'input' can be parsed as a Recipient User ID
func ValidateRecipientUserID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRecipientUserID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Recipient User ID
func (id RecipientUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/notifications/%s/recipientUsers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, string(id.NotificationName), id.UserId)
}

// Segments returns a slice of Resource ID Segments which comprise this Recipient User ID
func (id RecipientUserId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticNotifications", "notifications", "notifications"),
		resourceids.ConstantSegment("notificationName", PossibleValuesForNotificationName(), "AccountClosedPublisher"),
		resourceids.StaticSegment("staticRecipientUsers", "recipientUsers", "recipientUsers"),
		resourceids.UserSpecifiedSegment("userId", "userId"),
	}
}

// String returns a human-readable description of this Recipient User ID
func (id RecipientUserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Notification Name: %q", string(id.NotificationName)),
		fmt.Sprintf("User: %q", id.UserId),
	}
	return fmt.Sprintf("Recipient User (%s)", strings.Join(components, "\n"))
}
