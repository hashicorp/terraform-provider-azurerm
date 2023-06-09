package virtualnetworktap

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualNetworkTapId{}

// VirtualNetworkTapId is a struct representing the Resource ID for a Virtual Network Tap
type VirtualNetworkTapId struct {
	SubscriptionId        string
	ResourceGroupName     string
	VirtualNetworkTapName string
}

// NewVirtualNetworkTapID returns a new VirtualNetworkTapId struct
func NewVirtualNetworkTapID(subscriptionId string, resourceGroupName string, virtualNetworkTapName string) VirtualNetworkTapId {
	return VirtualNetworkTapId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		VirtualNetworkTapName: virtualNetworkTapName,
	}
}

// ParseVirtualNetworkTapID parses 'input' into a VirtualNetworkTapId
func ParseVirtualNetworkTapID(input string) (*VirtualNetworkTapId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkTapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkTapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualNetworkTapName, ok = parsed.Parsed["virtualNetworkTapName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkTapName", *parsed)
	}

	return &id, nil
}

// ParseVirtualNetworkTapIDInsensitively parses 'input' case-insensitively into a VirtualNetworkTapId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkTapIDInsensitively(input string) (*VirtualNetworkTapId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkTapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkTapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualNetworkTapName, ok = parsed.Parsed["virtualNetworkTapName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkTapName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualNetworkTapID checks that 'input' can be parsed as a Virtual Network Tap ID
func ValidateVirtualNetworkTapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkTapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Tap ID
func (id VirtualNetworkTapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkTaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkTapName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Tap ID
func (id VirtualNetworkTapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworkTaps", "virtualNetworkTaps", "virtualNetworkTaps"),
		resourceids.UserSpecifiedSegment("virtualNetworkTapName", "virtualNetworkTapValue"),
	}
}

// String returns a human-readable description of this Virtual Network Tap ID
func (id VirtualNetworkTapId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Tap Name: %q", id.VirtualNetworkTapName),
	}
	return fmt.Sprintf("Virtual Network Tap (%s)", strings.Join(components, "\n"))
}
