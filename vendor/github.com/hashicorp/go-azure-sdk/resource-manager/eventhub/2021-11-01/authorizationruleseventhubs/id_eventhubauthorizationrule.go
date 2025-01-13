package authorizationruleseventhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EventhubAuthorizationRuleId{})
}

var _ resourceids.ResourceId = &EventhubAuthorizationRuleId{}

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
	parser := resourceids.NewParserFromResourceIdType(&EventhubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventhubAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEventhubAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a EventhubAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseEventhubAuthorizationRuleIDInsensitively(input string) (*EventhubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventhubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventhubAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EventhubAuthorizationRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.EventhubName, ok = input.Parsed["eventhubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventhubName", input)
	}

	if id.AuthorizationRuleName, ok = input.Parsed["authorizationRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticEventhubs", "eventhubs", "eventhubs"),
		resourceids.UserSpecifiedSegment("eventhubName", "eventhubName"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleName"),
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
