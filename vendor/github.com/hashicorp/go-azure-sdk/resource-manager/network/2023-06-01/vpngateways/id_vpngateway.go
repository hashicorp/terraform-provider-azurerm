package vpngateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VpnGatewayId{}

// VpnGatewayId is a struct representing the Resource ID for a Vpn Gateway
type VpnGatewayId struct {
	SubscriptionId    string
	ResourceGroupName string
	VpnGatewayName    string
}

// NewVpnGatewayID returns a new VpnGatewayId struct
func NewVpnGatewayID(subscriptionId string, resourceGroupName string, vpnGatewayName string) VpnGatewayId {
	return VpnGatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VpnGatewayName:    vpnGatewayName,
	}
}

// ParseVpnGatewayID parses 'input' into a VpnGatewayId
func ParseVpnGatewayID(input string) (*VpnGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	return &id, nil
}

// ParseVpnGatewayIDInsensitively parses 'input' case-insensitively into a VpnGatewayId
// note: this method should only be used for API response data and not user input
func ParseVpnGatewayIDInsensitively(input string) (*VpnGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	return &id, nil
}

// ValidateVpnGatewayID checks that 'input' can be parsed as a Vpn Gateway ID
func ValidateVpnGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVpnGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vpn Gateway ID
func (id VpnGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vpn Gateway ID
func (id VpnGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnGateways", "vpnGateways", "vpnGateways"),
		resourceids.UserSpecifiedSegment("vpnGatewayName", "vpnGatewayValue"),
	}
}

// String returns a human-readable description of this Vpn Gateway ID
func (id VpnGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Gateway Name: %q", id.VpnGatewayName),
	}
	return fmt.Sprintf("Vpn Gateway (%s)", strings.Join(components, "\n"))
}
