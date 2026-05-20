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
	recaser.RegisterResourceId(&VirtualNetworkGatewayId{})
}

var _ resourceids.ResourceId = &VirtualNetworkGatewayId{}

// VirtualNetworkGatewayId is a struct representing the Resource ID for a Virtual Network Gateway
type VirtualNetworkGatewayId struct {
	SubscriptionId            string
	ResourceGroupName         string
	VirtualNetworkGatewayName string
}

// NewVirtualNetworkGatewayID returns a new VirtualNetworkGatewayId struct
func NewVirtualNetworkGatewayID(subscriptionId string, resourceGroupName string, virtualNetworkGatewayName string) VirtualNetworkGatewayId {
	return VirtualNetworkGatewayId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
	}
}

// ParseVirtualNetworkGatewayID parses 'input' into a VirtualNetworkGatewayId
func ParseVirtualNetworkGatewayID(input string) (*VirtualNetworkGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkGatewayIDInsensitively parses 'input' case-insensitively into a VirtualNetworkGatewayId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkGatewayIDInsensitively(input string) (*VirtualNetworkGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkGatewayId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateVirtualNetworkGatewayID checks that 'input' can be parsed as a Virtual Network Gateway ID
func ValidateVirtualNetworkGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Gateway ID
func (id VirtualNetworkGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Gateway ID
func (id VirtualNetworkGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworkGateways", "virtualNetworkGateways", "virtualNetworkGateways"),
		resourceids.UserSpecifiedSegment("virtualNetworkGatewayName", "virtualNetworkGatewayName"),
	}
}

// String returns a human-readable description of this Virtual Network Gateway ID
func (id VirtualNetworkGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Gateway Name: %q", id.VirtualNetworkGatewayName),
	}
	return fmt.Sprintf("Virtual Network Gateway (%s)", strings.Join(components, "\n"))
}
