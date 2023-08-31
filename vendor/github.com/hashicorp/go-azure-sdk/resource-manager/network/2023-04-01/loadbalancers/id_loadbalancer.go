package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LoadBalancerId{}

// LoadBalancerId is a struct representing the Resource ID for a Load Balancer
type LoadBalancerId struct {
	SubscriptionId    string
	ResourceGroupName string
	LoadBalancerName  string
}

// NewLoadBalancerID returns a new LoadBalancerId struct
func NewLoadBalancerID(subscriptionId string, resourceGroupName string, loadBalancerName string) LoadBalancerId {
	return LoadBalancerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LoadBalancerName:  loadBalancerName,
	}
}

// ParseLoadBalancerID parses 'input' into a LoadBalancerId
func ParseLoadBalancerID(input string) (*LoadBalancerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	return &id, nil
}

// ParseLoadBalancerIDInsensitively parses 'input' case-insensitively into a LoadBalancerId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancerIDInsensitively(input string) (*LoadBalancerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	return &id, nil
}

// ValidateLoadBalancerID checks that 'input' can be parsed as a Load Balancer ID
func ValidateLoadBalancerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoadBalancerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Load Balancer ID
func (id LoadBalancerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Load Balancer ID
func (id LoadBalancerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupValue"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
	}
}

// String returns a human-readable description of this Load Balancer ID
func (id LoadBalancerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
	}
	return fmt.Sprintf("Load Balancer (%s)", strings.Join(components, "\n"))
}
