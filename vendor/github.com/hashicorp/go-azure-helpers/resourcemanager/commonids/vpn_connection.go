// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VPNConnectionId{}

// VPNConnectionId is a struct representing the Resource ID for a V P N Connection
type VPNConnectionId struct {
	SubscriptionId    string
	ResourceGroupName string
	GatewayName       string
	ConnectionName    string
}

// NewVPNConnectionID returns a new VPNConnectionId struct
func NewVPNConnectionID(subscriptionId string, resourceGroupName string, gatewayName string, connectionName string) VPNConnectionId {
	return VPNConnectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GatewayName:       gatewayName,
		ConnectionName:    connectionName,
	}
}

// ParseVPNConnectionID parses 'input' into a VPNConnectionId
func ParseVPNConnectionID(input string) (*VPNConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VPNConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VPNConnectionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVPNConnectionIDInsensitively parses 'input' case-insensitively into a VPNConnectionId
// note: this method should only be used for API response data and not user input
func ParseVPNConnectionIDInsensitively(input string) (*VPNConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VPNConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VPNConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VPNConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GatewayName, ok = input.Parsed["gatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", input)
	}

	if id.ConnectionName, ok = input.Parsed["connectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionName", input)
	}

	return nil
}

// ValidateVPNConnectionID checks that 'input' can be parsed as a V P N Connection ID
func ValidateVPNConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVPNConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted V P N Connection ID
func (id VPNConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GatewayName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this V P N Connection ID
func (id VPNConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("vpnGateways", "vpnGateways", "vpnGateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayValue"),
		resourceids.StaticSegment("vpnConnections", "vpnConnections", "vpnConnections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionValue"),
	}
}

// String returns a human-readable description of this V P N Connection ID
func (id VPNConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("VPN Connection (%s)", strings.Join(components, "\n"))
}
