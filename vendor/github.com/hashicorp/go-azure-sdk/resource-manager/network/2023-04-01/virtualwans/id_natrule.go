package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NatRuleId{}

// NatRuleId is a struct representing the Resource ID for a Nat Rule
type NatRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	VpnGatewayName    string
	NatRuleName       string
}

// NewNatRuleID returns a new NatRuleId struct
func NewNatRuleID(subscriptionId string, resourceGroupName string, vpnGatewayName string, natRuleName string) NatRuleId {
	return NatRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VpnGatewayName:    vpnGatewayName,
		NatRuleName:       natRuleName,
	}
}

// ParseNatRuleID parses 'input' into a NatRuleId
func ParseNatRuleID(input string) (*NatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(NatRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NatRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	if id.NatRuleName, ok = parsed.Parsed["natRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "natRuleName", *parsed)
	}

	return &id, nil
}

// ParseNatRuleIDInsensitively parses 'input' case-insensitively into a NatRuleId
// note: this method should only be used for API response data and not user input
func ParseNatRuleIDInsensitively(input string) (*NatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(NatRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NatRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	if id.NatRuleName, ok = parsed.Parsed["natRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "natRuleName", *parsed)
	}

	return &id, nil
}

// ValidateNatRuleID checks that 'input' can be parsed as a Nat Rule ID
func ValidateNatRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNatRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Nat Rule ID
func (id NatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/natRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnGatewayName, id.NatRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Nat Rule ID
func (id NatRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnGateways", "vpnGateways", "vpnGateways"),
		resourceids.UserSpecifiedSegment("vpnGatewayName", "vpnGatewayValue"),
		resourceids.StaticSegment("staticNatRules", "natRules", "natRules"),
		resourceids.UserSpecifiedSegment("natRuleName", "natRuleValue"),
	}
}

// String returns a human-readable description of this Nat Rule ID
func (id NatRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Gateway Name: %q", id.VpnGatewayName),
		fmt.Sprintf("Nat Rule Name: %q", id.NatRuleName),
	}
	return fmt.Sprintf("Nat Rule (%s)", strings.Join(components, "\n"))
}
