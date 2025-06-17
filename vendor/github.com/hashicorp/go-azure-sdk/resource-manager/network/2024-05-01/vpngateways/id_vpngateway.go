package vpngateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VpnGatewayId{})
}

var _ resourceids.ResourceId = &VpnGatewayId{}

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
	parser := resourceids.NewParserFromResourceIdType(&VpnGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVpnGatewayIDInsensitively parses 'input' case-insensitively into a VpnGatewayId
// note: this method should only be used for API response data and not user input
func ParseVpnGatewayIDInsensitively(input string) (*VpnGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VpnGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VpnGatewayId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VpnGatewayName, ok = input.Parsed["vpnGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vpnGatewayName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vpnGatewayName", "vpnGatewayName"),
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
