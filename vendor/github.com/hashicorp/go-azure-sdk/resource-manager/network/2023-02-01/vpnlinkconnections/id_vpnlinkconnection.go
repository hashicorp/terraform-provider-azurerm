package vpnlinkconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VpnLinkConnectionId{}

// VpnLinkConnectionId is a struct representing the Resource ID for a Vpn Link Connection
type VpnLinkConnectionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	VpnGatewayName        string
	VpnConnectionName     string
	VpnLinkConnectionName string
}

// NewVpnLinkConnectionID returns a new VpnLinkConnectionId struct
func NewVpnLinkConnectionID(subscriptionId string, resourceGroupName string, vpnGatewayName string, vpnConnectionName string, vpnLinkConnectionName string) VpnLinkConnectionId {
	return VpnLinkConnectionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		VpnGatewayName:        vpnGatewayName,
		VpnConnectionName:     vpnConnectionName,
		VpnLinkConnectionName: vpnLinkConnectionName,
	}
}

// ParseVpnLinkConnectionID parses 'input' into a VpnLinkConnectionId
func ParseVpnLinkConnectionID(input string) (*VpnLinkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnLinkConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnLinkConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	if id.VpnConnectionName, ok = parsed.Parsed["vpnConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnConnectionName", *parsed)
	}

	if id.VpnLinkConnectionName, ok = parsed.Parsed["vpnLinkConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnLinkConnectionName", *parsed)
	}

	return &id, nil
}

// ParseVpnLinkConnectionIDInsensitively parses 'input' case-insensitively into a VpnLinkConnectionId
// note: this method should only be used for API response data and not user input
func ParseVpnLinkConnectionIDInsensitively(input string) (*VpnLinkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnLinkConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnLinkConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnGatewayName, ok = parsed.Parsed["vpnGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", *parsed)
	}

	if id.VpnConnectionName, ok = parsed.Parsed["vpnConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnConnectionName", *parsed)
	}

	if id.VpnLinkConnectionName, ok = parsed.Parsed["vpnLinkConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnLinkConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateVpnLinkConnectionID checks that 'input' can be parsed as a Vpn Link Connection ID
func ValidateVpnLinkConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVpnLinkConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vpn Link Connection ID
func (id VpnLinkConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s/vpnLinkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnGatewayName, id.VpnConnectionName, id.VpnLinkConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vpn Link Connection ID
func (id VpnLinkConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnGateways", "vpnGateways", "vpnGateways"),
		resourceids.UserSpecifiedSegment("vpnGatewayName", "vpnGatewayValue"),
		resourceids.StaticSegment("staticVpnConnections", "vpnConnections", "vpnConnections"),
		resourceids.UserSpecifiedSegment("vpnConnectionName", "vpnConnectionValue"),
		resourceids.StaticSegment("staticVpnLinkConnections", "vpnLinkConnections", "vpnLinkConnections"),
		resourceids.UserSpecifiedSegment("vpnLinkConnectionName", "vpnLinkConnectionValue"),
	}
}

// String returns a human-readable description of this Vpn Link Connection ID
func (id VpnLinkConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Gateway Name: %q", id.VpnGatewayName),
		fmt.Sprintf("Vpn Connection Name: %q", id.VpnConnectionName),
		fmt.Sprintf("Vpn Link Connection Name: %q", id.VpnLinkConnectionName),
	}
	return fmt.Sprintf("Vpn Link Connection (%s)", strings.Join(components, "\n"))
}
