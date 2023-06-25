package topicsauthorizationrule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TopicAuthorizationRuleId{}

// TopicAuthorizationRuleId is a struct representing the Resource ID for a Topic Authorization Rule
type TopicAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	TopicName             string
	AuthorizationRuleName string
}

// NewTopicAuthorizationRuleID returns a new TopicAuthorizationRuleId struct
func NewTopicAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string, authorizationRuleName string) TopicAuthorizationRuleId {
	return TopicAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		TopicName:             topicName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseTopicAuthorizationRuleID parses 'input' into a TopicAuthorizationRuleId
func ParseTopicAuthorizationRuleID(input string) (*TopicAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ParseTopicAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a TopicAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseTopicAuthorizationRuleIDInsensitively(input string) (*TopicAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(TopicAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TopicAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "topicName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ValidateTopicAuthorizationRuleID checks that 'input' can be parsed as a Topic Authorization Rule ID
func ValidateTopicAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTopicAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Topic Authorization Rule ID
func (id TopicAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Topic Authorization Rule ID
func (id TopicAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicValue"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Topic Authorization Rule ID
func (id TopicAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Topic Authorization Rule (%s)", strings.Join(components, "\n"))
}
