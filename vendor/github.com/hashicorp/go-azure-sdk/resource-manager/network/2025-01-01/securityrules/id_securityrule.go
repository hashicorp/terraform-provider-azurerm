package securityrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SecurityRuleId{})
}

var _ resourceids.ResourceId = &SecurityRuleId{}

// SecurityRuleId is a struct representing the Resource ID for a Security Rule
type SecurityRuleId struct {
	SubscriptionId           string
	ResourceGroupName        string
	NetworkSecurityGroupName string
	SecurityRuleName         string
}

// NewSecurityRuleID returns a new SecurityRuleId struct
func NewSecurityRuleID(subscriptionId string, resourceGroupName string, networkSecurityGroupName string, securityRuleName string) SecurityRuleId {
	return SecurityRuleId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		NetworkSecurityGroupName: networkSecurityGroupName,
		SecurityRuleName:         securityRuleName,
	}
}

// ParseSecurityRuleID parses 'input' into a SecurityRuleId
func ParseSecurityRuleID(input string) (*SecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecurityRuleIDInsensitively parses 'input' case-insensitively into a SecurityRuleId
// note: this method should only be used for API response data and not user input
func ParseSecurityRuleIDInsensitively(input string) (*SecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecurityRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecurityRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecurityRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkSecurityGroupName, ok = input.Parsed["networkSecurityGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityGroupName", input)
	}

	if id.SecurityRuleName, ok = input.Parsed["securityRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "securityRuleName", input)
	}

	return nil
}

// ValidateSecurityRuleID checks that 'input' can be parsed as a Security Rule ID
func ValidateSecurityRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Rule ID
func (id SecurityRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s/securityRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityGroupName, id.SecurityRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Rule ID
func (id SecurityRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityGroups", "networkSecurityGroups", "networkSecurityGroups"),
		resourceids.UserSpecifiedSegment("networkSecurityGroupName", "networkSecurityGroupName"),
		resourceids.StaticSegment("staticSecurityRules", "securityRules", "securityRules"),
		resourceids.UserSpecifiedSegment("securityRuleName", "securityRuleName"),
	}
}

// String returns a human-readable description of this Security Rule ID
func (id SecurityRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Group Name: %q", id.NetworkSecurityGroupName),
		fmt.Sprintf("Security Rule Name: %q", id.SecurityRuleName),
	}
	return fmt.Sprintf("Security Rule (%s)", strings.Join(components, "\n"))
}
