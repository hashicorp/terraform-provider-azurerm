package authorizationrulesnamespaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthorizationRuleId{})
}

var _ resourceids.ResourceId = &AuthorizationRuleId{}

// AuthorizationRuleId is a struct representing the Resource ID for a Authorization Rule
type AuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	AuthorizationRuleName string
}

// NewAuthorizationRuleID returns a new AuthorizationRuleId struct
func NewAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, authorizationRuleName string) AuthorizationRuleId {
	return AuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseAuthorizationRuleID parses 'input' into a AuthorizationRuleId
func ParseAuthorizationRuleID(input string) (*AuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a AuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseAuthorizationRuleIDInsensitively(input string) (*AuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthorizationRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AuthorizationRuleName, ok = input.Parsed["authorizationRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", input)
	}

	return nil
}

// ValidateAuthorizationRuleID checks that 'input' can be parsed as a Authorization Rule ID
func ValidateAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authorization Rule ID
func (id AuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Authorization Rule ID
func (id AuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleName"),
	}
}

// String returns a human-readable description of this Authorization Rule ID
func (id AuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Authorization Rule (%s)", strings.Join(components, "\n"))
}
