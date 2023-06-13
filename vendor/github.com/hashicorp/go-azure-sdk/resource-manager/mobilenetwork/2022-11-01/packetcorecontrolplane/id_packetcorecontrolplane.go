package packetcorecontrolplane

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PacketCoreControlPlaneId{}

// PacketCoreControlPlaneId is a struct representing the Resource ID for a Packet Core Control Plane
type PacketCoreControlPlaneId struct {
	SubscriptionId             string
	ResourceGroupName          string
	PacketCoreControlPlaneName string
}

// NewPacketCoreControlPlaneID returns a new PacketCoreControlPlaneId struct
func NewPacketCoreControlPlaneID(subscriptionId string, resourceGroupName string, packetCoreControlPlaneName string) PacketCoreControlPlaneId {
	return PacketCoreControlPlaneId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		PacketCoreControlPlaneName: packetCoreControlPlaneName,
	}
}

// ParsePacketCoreControlPlaneID parses 'input' into a PacketCoreControlPlaneId
func ParsePacketCoreControlPlaneID(input string) (*PacketCoreControlPlaneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PacketCoreControlPlaneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PacketCoreControlPlaneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PacketCoreControlPlaneName, ok = parsed.Parsed["packetCoreControlPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreControlPlaneName", *parsed)
	}

	return &id, nil
}

// ParsePacketCoreControlPlaneIDInsensitively parses 'input' case-insensitively into a PacketCoreControlPlaneId
// note: this method should only be used for API response data and not user input
func ParsePacketCoreControlPlaneIDInsensitively(input string) (*PacketCoreControlPlaneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PacketCoreControlPlaneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PacketCoreControlPlaneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PacketCoreControlPlaneName, ok = parsed.Parsed["packetCoreControlPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreControlPlaneName", *parsed)
	}

	return &id, nil
}

// ValidatePacketCoreControlPlaneID checks that 'input' can be parsed as a Packet Core Control Plane ID
func ValidatePacketCoreControlPlaneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePacketCoreControlPlaneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Packet Core Control Plane ID
func (id PacketCoreControlPlaneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Packet Core Control Plane ID
func (id PacketCoreControlPlaneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticPacketCoreControlPlanes", "packetCoreControlPlanes", "packetCoreControlPlanes"),
		resourceids.UserSpecifiedSegment("packetCoreControlPlaneName", "packetCoreControlPlaneValue"),
	}
}

// String returns a human-readable description of this Packet Core Control Plane ID
func (id PacketCoreControlPlaneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Packet Core Control Plane Name: %q", id.PacketCoreControlPlaneName),
	}
	return fmt.Sprintf("Packet Core Control Plane (%s)", strings.Join(components, "\n"))
}
