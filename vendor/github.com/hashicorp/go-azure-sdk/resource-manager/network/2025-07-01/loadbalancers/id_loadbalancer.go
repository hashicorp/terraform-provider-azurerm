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
	recaser.RegisterResourceId(&LoadBalancerId{})
}

var _ resourceids.ResourceId = &LoadBalancerId{}

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
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLoadBalancerIDInsensitively parses 'input' case-insensitively into a LoadBalancerId
// note: this method should only be used for API response data and not user input
func ParseLoadBalancerIDInsensitively(input string) (*LoadBalancerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoadBalancerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadBalancerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LoadBalancerId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
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
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupName"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerName"),
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
