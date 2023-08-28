package loadbalancers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProbeId{}

// ProbeId is a struct representing the Resource ID for a Probe
type ProbeId struct {
	SubscriptionId    string
	ResourceGroupName string
	LoadBalancerName  string
	ProbeName         string
}

// NewProbeID returns a new ProbeId struct
func NewProbeID(subscriptionId string, resourceGroupName string, loadBalancerName string, probeName string) ProbeId {
	return ProbeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LoadBalancerName:  loadBalancerName,
		ProbeName:         probeName,
	}
}

// ParseProbeID parses 'input' into a ProbeId
func ParseProbeID(input string) (*ProbeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProbeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProbeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.ProbeName, ok = parsed.Parsed["probeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "probeName", *parsed)
	}

	return &id, nil
}

// ParseProbeIDInsensitively parses 'input' case-insensitively into a ProbeId
// note: this method should only be used for API response data and not user input
func ParseProbeIDInsensitively(input string) (*ProbeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProbeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProbeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadBalancerName, ok = parsed.Parsed["loadBalancerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadBalancerName", *parsed)
	}

	if id.ProbeName, ok = parsed.Parsed["probeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "probeName", *parsed)
	}

	return &id, nil
}

// ValidateProbeID checks that 'input' can be parsed as a Probe ID
func ValidateProbeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProbeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Probe ID
func (id ProbeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/probes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName, id.ProbeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Probe ID
func (id ProbeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLoadBalancers", "loadBalancers", "loadBalancers"),
		resourceids.UserSpecifiedSegment("loadBalancerName", "loadBalancerValue"),
		resourceids.StaticSegment("staticProbes", "probes", "probes"),
		resourceids.UserSpecifiedSegment("probeName", "probeValue"),
	}
}

// String returns a human-readable description of this Probe ID
func (id ProbeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Balancer Name: %q", id.LoadBalancerName),
		fmt.Sprintf("Probe Name: %q", id.ProbeName),
	}
	return fmt.Sprintf("Probe (%s)", strings.Join(components, "\n"))
}
