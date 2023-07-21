package attacheddatanetwork

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AttachedDataNetworkId{}

// AttachedDataNetworkId is a struct representing the Resource ID for a Attached Data Network
type AttachedDataNetworkId struct {
	SubscriptionId             string
	ResourceGroupName          string
	PacketCoreControlPlaneName string
	PacketCoreDataPlaneName    string
	AttachedDataNetworkName    string
}

// NewAttachedDataNetworkID returns a new AttachedDataNetworkId struct
func NewAttachedDataNetworkID(subscriptionId string, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, attachedDataNetworkName string) AttachedDataNetworkId {
	return AttachedDataNetworkId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		PacketCoreControlPlaneName: packetCoreControlPlaneName,
		PacketCoreDataPlaneName:    packetCoreDataPlaneName,
		AttachedDataNetworkName:    attachedDataNetworkName,
	}
}

// ParseAttachedDataNetworkID parses 'input' into a AttachedDataNetworkId
func ParseAttachedDataNetworkID(input string) (*AttachedDataNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(AttachedDataNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AttachedDataNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PacketCoreControlPlaneName, ok = parsed.Parsed["packetCoreControlPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreControlPlaneName", *parsed)
	}

	if id.PacketCoreDataPlaneName, ok = parsed.Parsed["packetCoreDataPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreDataPlaneName", *parsed)
	}

	if id.AttachedDataNetworkName, ok = parsed.Parsed["attachedDataNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "attachedDataNetworkName", *parsed)
	}

	return &id, nil
}

// ParseAttachedDataNetworkIDInsensitively parses 'input' case-insensitively into a AttachedDataNetworkId
// note: this method should only be used for API response data and not user input
func ParseAttachedDataNetworkIDInsensitively(input string) (*AttachedDataNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(AttachedDataNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AttachedDataNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PacketCoreControlPlaneName, ok = parsed.Parsed["packetCoreControlPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreControlPlaneName", *parsed)
	}

	if id.PacketCoreDataPlaneName, ok = parsed.Parsed["packetCoreDataPlaneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCoreDataPlaneName", *parsed)
	}

	if id.AttachedDataNetworkName, ok = parsed.Parsed["attachedDataNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "attachedDataNetworkName", *parsed)
	}

	return &id, nil
}

// ValidateAttachedDataNetworkID checks that 'input' can be parsed as a Attached Data Network ID
func ValidateAttachedDataNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAttachedDataNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Attached Data Network ID
func (id AttachedDataNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/%s/packetCoreDataPlanes/%s/attachedDataNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PacketCoreControlPlaneName, id.PacketCoreDataPlaneName, id.AttachedDataNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Attached Data Network ID
func (id AttachedDataNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticPacketCoreControlPlanes", "packetCoreControlPlanes", "packetCoreControlPlanes"),
		resourceids.UserSpecifiedSegment("packetCoreControlPlaneName", "packetCoreControlPlaneValue"),
		resourceids.StaticSegment("staticPacketCoreDataPlanes", "packetCoreDataPlanes", "packetCoreDataPlanes"),
		resourceids.UserSpecifiedSegment("packetCoreDataPlaneName", "packetCoreDataPlaneValue"),
		resourceids.StaticSegment("staticAttachedDataNetworks", "attachedDataNetworks", "attachedDataNetworks"),
		resourceids.UserSpecifiedSegment("attachedDataNetworkName", "attachedDataNetworkValue"),
	}
}

// String returns a human-readable description of this Attached Data Network ID
func (id AttachedDataNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Packet Core Control Plane Name: %q", id.PacketCoreControlPlaneName),
		fmt.Sprintf("Packet Core Data Plane Name: %q", id.PacketCoreDataPlaneName),
		fmt.Sprintf("Attached Data Network Name: %q", id.AttachedDataNetworkName),
	}
	return fmt.Sprintf("Attached Data Network (%s)", strings.Join(components, "\n"))
}
