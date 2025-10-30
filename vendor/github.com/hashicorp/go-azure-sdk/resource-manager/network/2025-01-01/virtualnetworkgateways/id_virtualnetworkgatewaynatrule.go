package virtualnetworkgateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkGatewayNatRuleId{})
}

var _ resourceids.ResourceId = &VirtualNetworkGatewayNatRuleId{}

// VirtualNetworkGatewayNatRuleId is a struct representing the Resource ID for a Virtual Network Gateway Nat Rule
type VirtualNetworkGatewayNatRuleId struct {
	SubscriptionId            string
	ResourceGroupName         string
	VirtualNetworkGatewayName string
	NatRuleName               string
}

// NewVirtualNetworkGatewayNatRuleID returns a new VirtualNetworkGatewayNatRuleId struct
func NewVirtualNetworkGatewayNatRuleID(subscriptionId string, resourceGroupName string, virtualNetworkGatewayName string, natRuleName string) VirtualNetworkGatewayNatRuleId {
	return VirtualNetworkGatewayNatRuleId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
		NatRuleName:               natRuleName,
	}
}

// ParseVirtualNetworkGatewayNatRuleID parses 'input' into a VirtualNetworkGatewayNatRuleId
func ParseVirtualNetworkGatewayNatRuleID(input string) (*VirtualNetworkGatewayNatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkGatewayNatRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkGatewayNatRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkGatewayNatRuleIDInsensitively parses 'input' case-insensitively into a VirtualNetworkGatewayNatRuleId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkGatewayNatRuleIDInsensitively(input string) (*VirtualNetworkGatewayNatRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkGatewayNatRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkGatewayNatRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkGatewayNatRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualNetworkGatewayName, ok = input.Parsed["virtualNetworkGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkGatewayName", input)
	}

	if id.NatRuleName, ok = input.Parsed["natRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "natRuleName", input)
	}

	return nil
}

// ValidateVirtualNetworkGatewayNatRuleID checks that 'input' can be parsed as a Virtual Network Gateway Nat Rule ID
func ValidateVirtualNetworkGatewayNatRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkGatewayNatRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Gateway Nat Rule ID
func (id VirtualNetworkGatewayNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s/natRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName, id.NatRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Gateway Nat Rule ID
func (id VirtualNetworkGatewayNatRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworkGateways", "virtualNetworkGateways", "virtualNetworkGateways"),
		resourceids.UserSpecifiedSegment("virtualNetworkGatewayName", "virtualNetworkGatewayName"),
		resourceids.StaticSegment("staticNatRules", "natRules", "natRules"),
		resourceids.UserSpecifiedSegment("natRuleName", "natRuleName"),
	}
}

// String returns a human-readable description of this Virtual Network Gateway Nat Rule ID
func (id VirtualNetworkGatewayNatRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Gateway Name: %q", id.VirtualNetworkGatewayName),
		fmt.Sprintf("Nat Rule Name: %q", id.NatRuleName),
	}
	return fmt.Sprintf("Virtual Network Gateway Nat Rule (%s)", strings.Join(components, "\n"))
}
