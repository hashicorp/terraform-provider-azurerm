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
	recaser.RegisterResourceId(&LoadBalancerBackendAddressPoolId{})
}

var _ resourceids.ResourceId = &LoadBalancerBackendAddressPoolId{}

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
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerBackendAddressPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerBackendAddressPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLoadBalancerBackendAddressPoolIDInsensitively parses 'input' case-insensitively into a LoadBalancerBackendAddressPoolId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancerBackendAddressPoolIDInsensitively(input string) (*LoadBalancerBackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerBackendAddressPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerBackendAddressPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LoadBalancerBackendAddressPoolId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackendAddressPoolName, ok = input.Parsed["backendAddressPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backendAddressPoolName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerName"),
		resourceids.StaticSegment("staticBackendAddressPools", "backendAddressPools", "backendAddressPools"),
		resourceids.UserSpecifiedSegment("backendAddressPoolName", "backendAddressPoolName"),
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
