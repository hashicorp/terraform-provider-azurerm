package simpolicy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SimPolicyId{}

// SimPolicyId is a struct representing the Resource ID for a Sim Policy
type SimPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	MobileNetworkName string
	SimPolicyName     string
}

// NewSimPolicyID returns a new SimPolicyId struct
func NewSimPolicyID(subscriptionId string, resourceGroupName string, mobileNetworkName string, simPolicyName string) SimPolicyId {
	return SimPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MobileNetworkName: mobileNetworkName,
		SimPolicyName:     simPolicyName,
	}
}

// ParseSimPolicyID parses 'input' into a SimPolicyId
func ParseSimPolicyID(input string) (*SimPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(SimPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SimPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SimPolicyName, ok = parsed.Parsed["simPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "simPolicyName", *parsed)
	}

	return &id, nil
}

// ParseSimPolicyIDInsensitively parses 'input' case-insensitively into a SimPolicyId
// note: this method should only be used for API response data and not user input
func ParseSimPolicyIDInsensitively(input string) (*SimPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(SimPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SimPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SimPolicyName, ok = parsed.Parsed["simPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "simPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateSimPolicyID checks that 'input' can be parsed as a Sim Policy ID
func ValidateSimPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSimPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sim Policy ID
func (id SimPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/mobileNetworks/%s/simPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName, id.SimPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sim Policy ID
func (id SimPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticMobileNetworks", "mobileNetworks", "mobileNetworks"),
		resourceids.UserSpecifiedSegment("mobileNetworkName", "mobileNetworkValue"),
		resourceids.StaticSegment("staticSimPolicies", "simPolicies", "simPolicies"),
		resourceids.UserSpecifiedSegment("simPolicyName", "simPolicyValue"),
	}
}

// String returns a human-readable description of this Sim Policy ID
func (id SimPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mobile Network Name: %q", id.MobileNetworkName),
		fmt.Sprintf("Sim Policy Name: %q", id.SimPolicyName),
	}
	return fmt.Sprintf("Sim Policy (%s)", strings.Join(components, "\n"))
}
