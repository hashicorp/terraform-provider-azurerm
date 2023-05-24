// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualNetworkId{}

// VirtualNetworkId is a struct representing the Resource ID for a Virtual Network
type VirtualNetworkId struct {
	SubscriptionId     string
	ResourceGroupName  string
	VirtualNetworkName string
}

// NewVirtualNetworkID returns a new VirtualNetworkId struct
func NewVirtualNetworkID(subscriptionId string, resourceGroupName string, virtualNetworkName string) VirtualNetworkId {
	return VirtualNetworkId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		VirtualNetworkName: virtualNetworkName,
	}
}

// ParseVirtualNetworkID parses 'input' into a VirtualNetworkId
func ParseVirtualNetworkID(input string) (*VirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualNetworkName, ok = parsed.Parsed["virtualNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkName", *parsed)
	}

	return &id, nil
}

// ParseVirtualNetworkIDInsensitively parses 'input' case-insensitively into a VirtualNetworkId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkIDInsensitively(input string) (*VirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualNetworkName, ok = parsed.Parsed["virtualNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualNetworkID checks that 'input' can be parsed as a Virtual Network ID
func ValidateVirtualNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network ID
func (id VirtualNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network ID
func (id VirtualNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("virtualNetworks", "virtualNetworks", "virtualNetworks"),
		resourceids.UserSpecifiedSegment("virtualNetworkName", "virtualNetworksValue"),
	}
}

// String returns a human-readable description of this Virtual Network ID
func (id VirtualNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Name: %q", id.VirtualNetworkName),
	}
	return fmt.Sprintf("Virtual Network (%s)", strings.Join(components, "\n"))
}
