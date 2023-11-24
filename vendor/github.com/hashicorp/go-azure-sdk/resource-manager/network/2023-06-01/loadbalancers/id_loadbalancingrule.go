package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LoadBalancingRuleId{}

// LoadBalancingRuleId is a struct representing the Resource ID for a Load Balancing Rule
type LoadBalancingRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	LoadBalancerName      string
	LoadBalancingRuleName string
}

// NewLoadBalancingRuleID returns a new LoadBalancingRuleId struct
func NewLoadBalancingRuleID(subscriptionId string, resourceGroupName string, loadBalancerName string, loadBalancingRuleName string) LoadBalancingRuleId {
	return LoadBalancingRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		LoadBalancerName:      loadBalancerName,
		LoadBalancingRuleName: loadBalancingRuleName,
	}
}

// ParseLoadBalancingRuleID parses 'input' into a LoadBalancingRuleId
func ParseLoadBalancingRuleID(input string) (*LoadBalancingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancingRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancingRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.LoadBalancingRuleName, ok = parsed.Parsed["loadBalancingRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancingRuleName", *parsed)
	}

	return &id, nil
}

// ParseLoadBalancingRuleIDInsensitively parses 'input' case-insensitively into a LoadBalancingRuleId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancingRuleIDInsensitively(input string) (*LoadBalancingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancingRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancingRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.LoadBalancingRuleName, ok = parsed.Parsed["loadBalancingRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancingRuleName", *parsed)
	}

	return &id, nil
}

// ValidateLoadBalancingRuleID checks that 'input' can be parsed as a Load Balancing Rule ID
func ValidateLoadBalancingRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoadBalancingRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Load Balancing Rule ID
func (id LoadBalancingRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/loadBalancingRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.LoadBalancingRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Load Balancing Rule ID
func (id LoadBalancingRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticLoadBalancingRules", "loadBalancingRules", "loadBalancingRules"),
		resourceids.UserSpecifiedSegment("loadBalancingRuleName", "loadBalancingRuleValue"),
	}
}

// String returns a human-readable description of this Load Balancing Rule ID
func (id LoadBalancingRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Load Balancing Rule Name: %q", id.LoadBalancingRuleName),
	}
	return fmt.Sprintf("Load Balancing Rule (%s)", strings.Join(components, "\n"))
}
