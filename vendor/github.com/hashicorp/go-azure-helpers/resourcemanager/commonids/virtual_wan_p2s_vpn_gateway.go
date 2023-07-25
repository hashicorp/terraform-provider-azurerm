// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualWANP2SVPNGatewayId{}

// VirtualWANP2SVPNGatewayId is a struct representing the Resource ID for a Virtual WAN P2S VPN Gateway
type VirtualWANP2SVPNGatewayId struct {
	SubscriptionId    string
	ResourceGroupName string
	GatewayName       string
}

// NewVirtualWANP2SVPNGatewayID returns a new VirtualWANP2SVPNGatewayId struct
func NewVirtualWANP2SVPNGatewayID(subscriptionId string, resourceGroupName string, gatewayName string) VirtualWANP2SVPNGatewayId {
	return VirtualWANP2SVPNGatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GatewayName:       gatewayName,
	}
}

// ParseVirtualWANP2SVPNGatewayID parses 'input' into a VirtualWANP2SVPNGatewayId
func ParseVirtualWANP2SVPNGatewayID(input string) (*VirtualWANP2SVPNGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualWANP2SVPNGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualWANP2SVPNGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ParseVirtualWANP2SVPNGatewayIDInsensitively parses 'input' case-insensitively into a VirtualWANP2SVPNGatewayId
// note: this method should only be used for API response data and not user input
func ParseVirtualWANP2SVPNGatewayIDInsensitively(input string) (*VirtualWANP2SVPNGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualWANP2SVPNGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualWANP2SVPNGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualWANP2SVPNGatewayID checks that 'input' can be parsed as a Virtual WAN P2S VPN Gateway ID
func ValidateVirtualWANP2SVPNGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualWANP2SVPNGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual WAN P2S VPN Gateway ID
func (id VirtualWANP2SVPNGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/p2sVpnGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual WAN P2S VPN Gateway ID
func (id VirtualWANP2SVPNGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("p2sVpnGateways", "p2sVpnGateways", "p2sVpnGateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayValue"),
	}
}

// String returns a human-readable description of this Virtual WAN P2S VPN Gateway ID
func (id VirtualWANP2SVPNGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
	}
	return fmt.Sprintf("Virtual WAN P2S VPN Gateway (%s)", strings.Join(components, "\n"))
}
