package securityrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DefaultSecurityRuleId{}

// DefaultSecurityRuleId is a struct representing the Resource ID for a Default Security Rule
type DefaultSecurityRuleId struct {
	SubscriptionId           string
	ResourceGroupName        string
	NetworkSecurityGroupName string
	DefaultSecurityRuleName  string
}

// NewDefaultSecurityRuleID returns a new DefaultSecurityRuleId struct
func NewDefaultSecurityRuleID(subscriptionId string, resourceGroupName string, networkSecurityGroupName string, defaultSecurityRuleName string) DefaultSecurityRuleId {
	return DefaultSecurityRuleId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		NetworkSecurityGroupName: networkSecurityGroupName,
		DefaultSecurityRuleName:  defaultSecurityRuleName,
	}
}

// ParseDefaultSecurityRuleID parses 'input' into a DefaultSecurityRuleId
func ParseDefaultSecurityRuleID(input string) (*DefaultSecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DefaultSecurityRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DefaultSecurityRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkSecurityGroupName, ok = parsed.Parsed["networkSecurityGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityGroupName", *parsed)
	}

	if id.DefaultSecurityRuleName, ok = parsed.Parsed["defaultSecurityRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "defaultSecurityRuleName", *parsed)
	}

	return &id, nil
}

// ParseDefaultSecurityRuleIDInsensitively parses 'input' case-insensitively into a DefaultSecurityRuleId
// note: this method should only be used for API response data and not user input
func ParseDefaultSecurityRuleIDInsensitively(input string) (*DefaultSecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DefaultSecurityRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DefaultSecurityRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkSecurityGroupName, ok = parsed.Parsed["networkSecurityGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityGroupName", *parsed)
	}

	if id.DefaultSecurityRuleName, ok = parsed.Parsed["defaultSecurityRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "defaultSecurityRuleName", *parsed)
	}

	return &id, nil
}

// ValidateDefaultSecurityRuleID checks that 'input' can be parsed as a Default Security Rule ID
func ValidateDefaultSecurityRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDefaultSecurityRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Default Security Rule ID
func (id DefaultSecurityRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s/defaultSecurityRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityGroupName, id.DefaultSecurityRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Default Security Rule ID
func (id DefaultSecurityRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityGroups", "networkSecurityGroups", "networkSecurityGroups"),
		resourceids.UserSpecifiedSegment("networkSecurityGroupName", "networkSecurityGroupValue"),
		resourceids.StaticSegment("staticDefaultSecurityRules", "defaultSecurityRules", "defaultSecurityRules"),
		resourceids.UserSpecifiedSegment("defaultSecurityRuleName", "defaultSecurityRuleValue"),
	}
}

// String returns a human-readable description of this Default Security Rule ID
func (id DefaultSecurityRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Group Name: %q", id.NetworkSecurityGroupName),
		fmt.Sprintf("Default Security Rule Name: %q", id.DefaultSecurityRuleName),
	}
	return fmt.Sprintf("Default Security Rule (%s)", strings.Join(components, "\n"))
}
