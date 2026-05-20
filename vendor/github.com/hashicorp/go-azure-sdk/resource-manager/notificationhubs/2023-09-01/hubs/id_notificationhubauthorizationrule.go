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
	recaser.RegisterResourceId(&NotificationHubAuthorizationRuleId{})
}

var _ resourceids.ResourceId = &NotificationHubAuthorizationRuleId{}

// NotificationHubAuthorizationRuleId is a struct representing the Resource ID for a Notification Hub Authorization Rule
type NotificationHubAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	NotificationHubName   string
	AuthorizationRuleName string
}

// NewNotificationHubAuthorizationRuleID returns a new NotificationHubAuthorizationRuleId struct
func NewNotificationHubAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, notificationHubName string, authorizationRuleName string) NotificationHubAuthorizationRuleId {
	return NotificationHubAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		NotificationHubName:   notificationHubName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseNotificationHubAuthorizationRuleID parses 'input' into a NotificationHubAuthorizationRuleId
func ParseNotificationHubAuthorizationRuleID(input string) (*NotificationHubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationHubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationHubAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNotificationHubAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a NotificationHubAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseNotificationHubAuthorizationRuleIDInsensitively(input string) (*NotificationHubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NotificationHubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NotificationHubAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NotificationHubAuthorizationRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AuthorizationRuleName, ok = input.Parsed["authorizationRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", input)
	}

	return nil
}

// ValidateNotificationHubAuthorizationRuleID checks that 'input' can be parsed as a Notification Hub Authorization Rule ID
func ValidateNotificationHubAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationHubAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleName"),
	}
}

// String returns a human-readable description of this Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Notification Hub Name: %q", id.NotificationHubName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Notification Hub Authorization Rule (%s)", strings.Join(components, "\n"))
}
