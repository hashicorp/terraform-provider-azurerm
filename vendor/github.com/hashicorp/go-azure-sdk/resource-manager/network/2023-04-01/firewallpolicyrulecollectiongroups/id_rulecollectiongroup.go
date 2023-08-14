package firewallpolicyrulecollectiongroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RuleCollectionGroupId{}

// RuleCollectionGroupId is a struct representing the Resource ID for a Rule Collection Group
type RuleCollectionGroupId struct {
	SubscriptionId          string
	ResourceGroupName       string
	FirewallPolicyName      string
	RuleCollectionGroupName string
}

// NewRuleCollectionGroupID returns a new RuleCollectionGroupId struct
func NewRuleCollectionGroupID(subscriptionId string, resourceGroupName string, firewallPolicyName string, ruleCollectionGroupName string) RuleCollectionGroupId {
	return RuleCollectionGroupId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		FirewallPolicyName:      firewallPolicyName,
		RuleCollectionGroupName: ruleCollectionGroupName,
	}
}

// ParseRuleCollectionGroupID parses 'input' into a RuleCollectionGroupId
func ParseRuleCollectionGroupID(input string) (*RuleCollectionGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleCollectionGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleCollectionGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FirewallPolicyName, ok = parsed.Parsed["firewallPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "firewallPolicyName", *parsed)
	}

	if id.RuleCollectionGroupName, ok = parsed.Parsed["ruleCollectionGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ruleCollectionGroupName", *parsed)
	}

	return &id, nil
}

// ParseRuleCollectionGroupIDInsensitively parses 'input' case-insensitively into a RuleCollectionGroupId
// note: this method should only be used for API response data and not user input
func ParseRuleCollectionGroupIDInsensitively(input string) (*RuleCollectionGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleCollectionGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleCollectionGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FirewallPolicyName, ok = parsed.Parsed["firewallPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "firewallPolicyName", *parsed)
	}

	if id.RuleCollectionGroupName, ok = parsed.Parsed["ruleCollectionGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ruleCollectionGroupName", *parsed)
	}

	return &id, nil
}

// ValidateRuleCollectionGroupID checks that 'input' can be parsed as a Rule Collection Group ID
func ValidateRuleCollectionGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleCollectionGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule Collection Group ID
func (id RuleCollectionGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s/ruleCollectionGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FirewallPolicyName, id.RuleCollectionGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule Collection Group ID
func (id RuleCollectionGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFirewallPolicies", "firewallPolicies", "firewallPolicies"),
		resourceids.UserSpecifiedSegment("firewallPolicyName", "firewallPolicyValue"),
		resourceids.StaticSegment("staticRuleCollectionGroups", "ruleCollectionGroups", "ruleCollectionGroups"),
		resourceids.UserSpecifiedSegment("ruleCollectionGroupName", "ruleCollectionGroupValue"),
	}
}

// String returns a human-readable description of this Rule Collection Group ID
func (id RuleCollectionGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Firewall Policy Name: %q", id.FirewallPolicyName),
		fmt.Sprintf("Rule Collection Group Name: %q", id.RuleCollectionGroupName),
	}
	return fmt.Sprintf("Rule Collection Group (%s)", strings.Join(components, "\n"))
}
