package authorizationruleseventhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = EventhubAuthorizationRuleId{}

// EventhubAuthorizationRuleId is a struct representing the Resource ID for a Eventhub Authorization Rule
type EventhubAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	EventhubName          string
	AuthorizationRuleName string
}

// NewEventhubAuthorizationRuleID returns a new EventhubAuthorizationRuleId struct
func NewEventhubAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, eventhubName string, authorizationRuleName string) EventhubAuthorizationRuleId {
	return EventhubAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		EventhubName:          eventhubName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseEventhubAuthorizationRuleID parses 'input' into a EventhubAuthorizationRuleId
func ParseEventhubAuthorizationRuleID(input string) (*EventhubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(EventhubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EventhubAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.EventhubName, ok = parsed.Parsed["eventhubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventhubName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ParseEventhubAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a EventhubAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseEventhubAuthorizationRuleIDInsensitively(input string) (*EventhubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(EventhubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EventhubAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.EventhubName, ok = parsed.Parsed["eventhubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventhubName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ValidateEventhubAuthorizationRuleID checks that 'input' can be parsed as a Eventhub Authorization Rule ID
func ValidateEventhubAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEventhubAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Eventhub Authorization Rule ID
func (id EventhubAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Eventhub Authorization Rule ID
func (id EventhubAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticEventhubs", "eventhubs", "eventhubs"),
		resourceids.UserSpecifiedSegment("eventhubName", "eventhubValue"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Eventhub Authorization Rule ID
func (id EventhubAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Eventhub Name: %q", id.EventhubName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Eventhub Authorization Rule (%s)", strings.Join(components, "\n"))
}
