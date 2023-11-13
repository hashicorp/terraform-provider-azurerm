package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BackendAddressPoolId{}

// BackendAddressPoolId is a struct representing the Resource ID for a Backend Address Pool
type BackendAddressPoolId struct {
	SubscriptionId         string
	ResourceGroupName      string
	LoadBalancerName       string
	BackendAddressPoolName string
}

// NewBackendAddressPoolID returns a new BackendAddressPoolId struct
func NewBackendAddressPoolID(subscriptionId string, resourceGroupName string, loadBalancerName string, backendAddressPoolName string) BackendAddressPoolId {
	return BackendAddressPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		LoadBalancerName:       loadBalancerName,
		BackendAddressPoolName: backendAddressPoolName,
	}
}

// ParseBackendAddressPoolID parses 'input' into a BackendAddressPoolId
func ParseBackendAddressPoolID(input string) (*BackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackendAddressPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackendAddressPoolId{}

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

// ParseBackendAddressPoolIDInsensitively parses 'input' case-insensitively into a BackendAddressPoolId
// note: this method should only be used for API response data and not user input
func ParseBackendAddressPoolIDInsensitively(input string) (*BackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackendAddressPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackendAddressPoolId{}

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

// ValidateBackendAddressPoolID checks that 'input' can be parsed as a Backend Address Pool ID
func ValidateBackendAddressPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackendAddressPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backend Address Pool ID
func (id BackendAddressPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.BackendAddressPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backend Address Pool ID
func (id BackendAddressPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupValue"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticBackendAddressPools", "backendAddressPools", "backendAddressPools"),
		resourceids.UserSpecifiedSegment("backendAddressPoolName", "backendAddressPoolValue"),
	}
}

// String returns a human-readable description of this Backend Address Pool ID
func (id BackendAddressPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Backend Address Pool Name: %q", id.BackendAddressPoolName),
	}
	return fmt.Sprintf("Backend Address Pool (%s)", strings.Join(components, "\n"))
}
