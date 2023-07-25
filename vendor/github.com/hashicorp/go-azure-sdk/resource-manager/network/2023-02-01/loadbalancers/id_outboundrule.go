package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = OutboundRuleId{}

// OutboundRuleId is a struct representing the Resource ID for a Outbound Rule
type OutboundRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	LoadBalancerName  string
	OutboundRuleName  string
}

// NewOutboundRuleID returns a new OutboundRuleId struct
func NewOutboundRuleID(subscriptionId string, resourceGroupName string, loadBalancerName string, outboundRuleName string) OutboundRuleId {
	return OutboundRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LoadBalancerName:  loadBalancerName,
		OutboundRuleName:  outboundRuleName,
	}
}

// ParseOutboundRuleID parses 'input' into a OutboundRuleId
func ParseOutboundRuleID(input string) (*OutboundRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutboundRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutboundRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.OutboundRuleName, ok = parsed.Parsed["outboundRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "outboundRuleName", *parsed)
	}

	return &id, nil
}

// ParseOutboundRuleIDInsensitively parses 'input' case-insensitively into a OutboundRuleId
// note: this method should only be used for API response data and not user input
func ParseOutboundRuleIDInsensitively(input string) (*OutboundRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutboundRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutboundRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.OutboundRuleName, ok = parsed.Parsed["outboundRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "outboundRuleName", *parsed)
	}

	return &id, nil
}

// ValidateOutboundRuleID checks that 'input' can be parsed as a Outbound Rule ID
func ValidateOutboundRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOutboundRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Outbound Rule ID
func (id OutboundRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/outboundRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.OutboundRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Outbound Rule ID
func (id OutboundRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticOutboundRules", "outboundRules", "outboundRules"),
		resourceids.UserSpecifiedSegment("outboundRuleName", "outboundRuleValue"),
	}
}

// String returns a human-readable description of this Outbound Rule ID
func (id OutboundRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Outbound Rule Name: %q", id.OutboundRuleName),
	}
	return fmt.Sprintf("Outbound Rule (%s)", strings.Join(components, "\n"))
}
