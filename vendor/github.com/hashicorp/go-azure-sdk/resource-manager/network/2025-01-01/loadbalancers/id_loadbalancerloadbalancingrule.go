package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LoadBalancerLoadBalancingRuleId{})
}

var _ resourceids.ResourceId = &LoadBalancerLoadBalancingRuleId{}

// LoadBalancerLoadBalancingRuleId is a struct representing the Resource ID for a Load Balancer Load Balancing Rule
type LoadBalancerLoadBalancingRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	LoadBalancerName      string
	LoadBalancingRuleName string
}

// NewLoadBalancerLoadBalancingRuleID returns a new LoadBalancerLoadBalancingRuleId struct
func NewLoadBalancerLoadBalancingRuleID(subscriptionId string, resourceGroupName string, loadBalancerName string, loadBalancingRuleName string) LoadBalancerLoadBalancingRuleId {
	return LoadBalancerLoadBalancingRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		LoadBalancerName:      loadBalancerName,
		LoadBalancingRuleName: loadBalancingRuleName,
	}
}

// ParseLoadBalancerLoadBalancingRuleID parses 'input' into a LoadBalancerLoadBalancingRuleId
func ParseLoadBalancerLoadBalancingRuleID(input string) (*LoadBalancerLoadBalancingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerLoadBalancingRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerLoadBalancingRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLoadBalancerLoadBalancingRuleIDInsensitively parses 'input' case-insensitively into a LoadBalancerLoadBalancingRuleId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancerLoadBalancingRuleIDInsensitively(input string) (*LoadBalancerLoadBalancingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerLoadBalancingRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerLoadBalancingRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LoadBalancerLoadBalancingRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LoadBalancerName, ok = input.Parsed["loadBalancerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", input)
	}

	if id.LoadBalancingRuleName, ok = input.Parsed["loadBalancingRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "loadBalancingRuleName", input)
	}

	return nil
}

// ValidateLoadBalancerLoadBalancingRuleID checks that 'input' can be parsed as a Load Balancer Load Balancing Rule ID
func ValidateLoadBalancerLoadBalancingRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoadBalancerLoadBalancingRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Load Balancer Load Balancing Rule ID
func (id LoadBalancerLoadBalancingRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/loadBalancingRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.LoadBalancingRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Load Balancer Load Balancing Rule ID
func (id LoadBalancerLoadBalancingRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerName"),
		resourceids.StaticSegment("staticLoadBalancingRules", "loadBalancingRules", "loadBalancingRules"),
		resourceids.UserSpecifiedSegment("loadBalancingRuleName", "loadBalancingRuleName"),
	}
}

// String returns a human-readable description of this Load Balancer Load Balancing Rule ID
func (id LoadBalancerLoadBalancingRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Load Balancing Rule Name: %q", id.LoadBalancingRuleName),
	}
	return fmt.Sprintf("Load Balancer Load Balancing Rule (%s)", strings.Join(components, "\n"))
}
