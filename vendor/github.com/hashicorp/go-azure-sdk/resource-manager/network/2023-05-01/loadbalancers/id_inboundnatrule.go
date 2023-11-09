package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = InboundNatRuleId{}

// InboundNatRuleId is a struct representing the Resource ID for a Inbound Nat Rule
type InboundNatRuleId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LoadBalancerName   string
	InboundNatRuleName string
}

// NewInboundNatRuleID returns a new InboundNatRuleId struct
func NewInboundNatRuleID(subscriptionId string, resourceGroupName string, loadBalancerName string, inboundNatRuleName string) InboundNatRuleId {
	return InboundNatRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LoadBalancerName:   loadBalancerName,
		InboundNatRuleName: inboundNatRuleName,
	}
}

// ParseInboundNatRuleID parses 'input' into a InboundNatRuleId
func ParseInboundNatRuleID(input string) (*InboundNatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundNatRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundNatRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.InboundNatRuleName, ok = parsed.Parsed["inboundNatRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundNatRuleName", *parsed)
	}

	return &id, nil
}

// ParseInboundNatRuleIDInsensitively parses 'input' case-insensitively into a InboundNatRuleId
// note: this method should only be used for API response data and not user input
func ParseInboundNatRuleIDInsensitively(input string) (*InboundNatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundNatRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundNatRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.InboundNatRuleName, ok = parsed.Parsed["inboundNatRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundNatRuleName", *parsed)
	}

	return &id, nil
}

// ValidateInboundNatRuleID checks that 'input' can be parsed as a Inbound Nat Rule ID
func ValidateInboundNatRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInboundNatRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Inbound Nat Rule ID
func (id InboundNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/inboundNatRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.InboundNatRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Inbound Nat Rule ID
func (id InboundNatRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticInboundNatRules", "inboundNatRules", "inboundNatRules"),
		resourceids.UserSpecifiedSegment("inboundNatRuleName", "inboundNatRuleValue"),
	}
}

// String returns a human-readable description of this Inbound Nat Rule ID
func (id InboundNatRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Inbound Nat Rule Name: %q", id.InboundNatRuleName),
	}
	return fmt.Sprintf("Inbound Nat Rule (%s)", strings.Join(components, "\n"))
}
