package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderLoadBalancerId{}

// ProviderLoadBalancerId is a struct representing the Resource ID for a Provider Load Balancer
type ProviderLoadBalancerId struct {
	SubscriptionId    string
	ResourceGroupName string
	LoadBalancerName  string
}

// NewProviderLoadBalancerID returns a new ProviderLoadBalancerId struct
func NewProviderLoadBalancerID(subscriptionId string, resourceGroupName string, loadBalancerName string) ProviderLoadBalancerId {
	return ProviderLoadBalancerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LoadBalancerName:  loadBalancerName,
	}
}

// ParseProviderLoadBalancerID parses 'input' into a ProviderLoadBalancerId
func ParseProviderLoadBalancerID(input string) (*ProviderLoadBalancerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLoadBalancerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLoadBalancerId{}

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

// ParseProviderLoadBalancerIDInsensitively parses 'input' case-insensitively into a ProviderLoadBalancerId
// note: this method should only be used for API response data and not user input
func ParseProviderLoadBalancerIDInsensitively(input string) (*ProviderLoadBalancerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLoadBalancerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLoadBalancerId{}

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

// ValidateProviderLoadBalancerID checks that 'input' can be parsed as a Provider Load Balancer ID
func ValidateProviderLoadBalancerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderLoadBalancerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Load Balancer ID
func (id ProviderLoadBalancerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Load Balancer ID
func (id ProviderLoadBalancerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
	}
}

// String returns a human-readable description of this Provider Load Balancer ID
func (id ProviderLoadBalancerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
	}
	return fmt.Sprintf("Provider Load Balancer (%s)", strings.Join(components, "\n"))
}
