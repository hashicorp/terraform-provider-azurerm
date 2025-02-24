package virtualnetworkpeerings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkPeeringId{})
}

var _ resourceids.ResourceId = &VirtualNetworkPeeringId{}

// VirtualNetworkPeeringId is a struct representing the Resource ID for a Virtual Network Peering
type VirtualNetworkPeeringId struct {
	SubscriptionId            string
	ResourceGroupName         string
	VirtualNetworkName        string
	VirtualNetworkPeeringName string
}

// NewVirtualNetworkPeeringID returns a new VirtualNetworkPeeringId struct
func NewVirtualNetworkPeeringID(subscriptionId string, resourceGroupName string, virtualNetworkName string, virtualNetworkPeeringName string) VirtualNetworkPeeringId {
	return VirtualNetworkPeeringId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		VirtualNetworkName:        virtualNetworkName,
		VirtualNetworkPeeringName: virtualNetworkPeeringName,
	}
}

// ParseVirtualNetworkPeeringID parses 'input' into a VirtualNetworkPeeringId
func ParseVirtualNetworkPeeringID(input string) (*VirtualNetworkPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkPeeringId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkPeeringId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkPeeringIDInsensitively parses 'input' case-insensitively into a VirtualNetworkPeeringId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkPeeringIDInsensitively(input string) (*VirtualNetworkPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkPeeringId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkPeeringId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkPeeringId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualNetworkName, ok = input.Parsed["virtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkName", input)
	}

	if id.VirtualNetworkPeeringName, ok = input.Parsed["virtualNetworkPeeringName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkPeeringName", input)
	}

	return nil
}

// ValidateVirtualNetworkPeeringID checks that 'input' can be parsed as a Virtual Network Peering ID
func ValidateVirtualNetworkPeeringID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkPeeringID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Peering ID
func (id VirtualNetworkPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/virtualNetworkPeerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, id.VirtualNetworkPeeringName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Peering ID
func (id VirtualNetworkPeeringId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworks", "virtualNetworks", "virtualNetworks"),
		resourceids.UserSpecifiedSegment("virtualNetworkName", "virtualNetworkName"),
		resourceids.StaticSegment("staticVirtualNetworkPeerings", "virtualNetworkPeerings", "virtualNetworkPeerings"),
		resourceids.UserSpecifiedSegment("virtualNetworkPeeringName", "virtualNetworkPeeringName"),
	}
}

// String returns a human-readable description of this Virtual Network Peering ID
func (id VirtualNetworkPeeringId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Name: %q", id.VirtualNetworkName),
		fmt.Sprintf("Virtual Network Peering Name: %q", id.VirtualNetworkPeeringName),
	}
	return fmt.Sprintf("Virtual Network Peering (%s)", strings.Join(components, "\n"))
}
