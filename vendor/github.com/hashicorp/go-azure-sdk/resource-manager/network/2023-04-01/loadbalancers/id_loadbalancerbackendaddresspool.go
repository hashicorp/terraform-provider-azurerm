package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LoadBalancerBackendAddressPoolId{}

// LoadBalancerBackendAddressPoolId is a struct representing the Resource ID for a Load Balancer Backend Address Pool
type LoadBalancerBackendAddressPoolId struct {
	SubscriptionId         string
	ResourceGroupName      string
	LoadBalancerName       string
	BackendAddressPoolName string
}

// NewLoadBalancerBackendAddressPoolID returns a new LoadBalancerBackendAddressPoolId struct
func NewLoadBalancerBackendAddressPoolID(subscriptionId string, resourceGroupName string, loadBalancerName string, backendAddressPoolName string) LoadBalancerBackendAddressPoolId {
	return LoadBalancerBackendAddressPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		LoadBalancerName:       loadBalancerName,
		BackendAddressPoolName: backendAddressPoolName,
	}
}

// ParseLoadBalancerBackendAddressPoolID parses 'input' into a LoadBalancerBackendAddressPoolId
func ParseLoadBalancerBackendAddressPoolID(input string) (*LoadBalancerBackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancerBackendAddressPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancerBackendAddressPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.BackendAddressPoolName, ok = parsed.Parsed["backendAddressPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backendAddressPoolName", *parsed)
	}

	return &id, nil
}

// ParseLoadBalancerBackendAddressPoolIDInsensitively parses 'input' case-insensitively into a LoadBalancerBackendAddressPoolId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancerBackendAddressPoolIDInsensitively(input string) (*LoadBalancerBackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadBalancerBackendAddressPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadBalancerBackendAddressPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.BackendAddressPoolName, ok = parsed.Parsed["backendAddressPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backendAddressPoolName", *parsed)
	}

	return &id, nil
}

// ValidateLoadBalancerBackendAddressPoolID checks that 'input' can be parsed as a Load Balancer Backend Address Pool ID
func ValidateLoadBalancerBackendAddressPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoadBalancerBackendAddressPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Load Balancer Backend Address Pool ID
func (id LoadBalancerBackendAddressPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.BackendAddressPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Load Balancer Backend Address Pool ID
func (id LoadBalancerBackendAddressPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticBackendAddressPools", "backendAddressPools", "backendAddressPools"),
		resourceids.UserSpecifiedSegment("backendAddressPoolName", "backendAddressPoolValue"),
	}
}

// String returns a human-readable description of this Load Balancer Backend Address Pool ID
func (id LoadBalancerBackendAddressPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Backend Address Pool Name: %q", id.BackendAddressPoolName),
	}
	return fmt.Sprintf("Load Balancer Backend Address Pool (%s)", strings.Join(components, "\n"))
}
