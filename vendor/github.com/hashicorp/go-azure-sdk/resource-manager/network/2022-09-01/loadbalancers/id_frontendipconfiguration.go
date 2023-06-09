package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FrontendIPConfigurationId{}

// FrontendIPConfigurationId is a struct representing the Resource ID for a Frontend I P Configuration
type FrontendIPConfigurationId struct {
	SubscriptionId              string
	ResourceGroupName           string
	LoadBalancerName            string
	FrontendIPConfigurationName string
}

// NewFrontendIPConfigurationID returns a new FrontendIPConfigurationId struct
func NewFrontendIPConfigurationID(subscriptionId string, resourceGroupName string, loadBalancerName string, frontendIPConfigurationName string) FrontendIPConfigurationId {
	return FrontendIPConfigurationId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		LoadBalancerName:            loadBalancerName,
		FrontendIPConfigurationName: frontendIPConfigurationName,
	}
}

// ParseFrontendIPConfigurationID parses 'input' into a FrontendIPConfigurationId
func ParseFrontendIPConfigurationID(input string) (*FrontendIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendIPConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendIPConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.FrontendIPConfigurationName, ok = parsed.Parsed["frontendIPConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendIPConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseFrontendIPConfigurationIDInsensitively parses 'input' case-insensitively into a FrontendIPConfigurationId
// note: this method should only be used for API response data and not user input
func ParseFrontendIPConfigurationIDInsensitively(input string) (*FrontendIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendIPConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendIPConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.FrontendIPConfigurationName, ok = parsed.Parsed["frontendIPConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendIPConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateFrontendIPConfigurationID checks that 'input' can be parsed as a Frontend I P Configuration ID
func ValidateFrontendIPConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFrontendIPConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Frontend I P Configuration ID
func (id FrontendIPConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.FrontendIPConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Frontend I P Configuration ID
func (id FrontendIPConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticFrontendIPConfigurations", "frontendIPConfigurations", "frontendIPConfigurations"),
		resourceids.UserSpecifiedSegment("frontendIPConfigurationName", "frontendIPConfigurationValue"),
	}
}

// String returns a human-readable description of this Frontend I P Configuration ID
func (id FrontendIPConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Frontend I P Configuration Name: %q", id.FrontendIPConfigurationName),
	}
	return fmt.Sprintf("Frontend I P Configuration (%s)", strings.Join(components, "\n"))
}
