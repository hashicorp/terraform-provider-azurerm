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
	recaser.RegisterResourceId(&BackendAddressPoolId{})
}

var _ resourceids.ResourceId = &BackendAddressPoolId{}

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
	parser := resourceids.NewParserFromResourceIdType(&BackendAddressPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackendAddressPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackendAddressPoolIDInsensitively parses 'input' case-insensitively into a BackendAddressPoolId
// note: this method should only be used for API response data and not user input
func ParseBackendAddressPoolIDInsensitively(input string) (*BackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackendAddressPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackendAddressPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackendAddressPoolId) FromParseResult(input resourceids.ParseResult) error {
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
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupName"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerName"),
		resourceids.StaticSegment("staticBackendAddressPools", "backendAddressPools", "backendAddressPools"),
		resourceids.UserSpecifiedSegment("backendAddressPoolName", "backendAddressPoolName"),
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
