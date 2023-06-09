package packetcaptures

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PacketCaptureId{}

// PacketCaptureId is a struct representing the Resource ID for a Packet Capture
type PacketCaptureId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkWatcherName string
	PacketCaptureName  string
}

// NewPacketCaptureID returns a new PacketCaptureId struct
func NewPacketCaptureID(subscriptionId string, resourceGroupName string, networkWatcherName string, packetCaptureName string) PacketCaptureId {
	return PacketCaptureId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkWatcherName: networkWatcherName,
		PacketCaptureName:  packetCaptureName,
	}
}

// ParsePacketCaptureID parses 'input' into a PacketCaptureId
func ParsePacketCaptureID(input string) (*PacketCaptureId, error) {
	parser := resourceids.NewParserFromResourceIdType(PacketCaptureId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PacketCaptureId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkWatcherName, ok = parsed.Parsed["networkWatcherName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkWatcherName", *parsed)
	}

	if id.PacketCaptureName, ok = parsed.Parsed["packetCaptureName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCaptureName", *parsed)
	}

	return &id, nil
}

// ParsePacketCaptureIDInsensitively parses 'input' case-insensitively into a PacketCaptureId
// note: this method should only be used for API response data and not user input
func ParsePacketCaptureIDInsensitively(input string) (*PacketCaptureId, error) {
	parser := resourceids.NewParserFromResourceIdType(PacketCaptureId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PacketCaptureId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkWatcherName, ok = parsed.Parsed["networkWatcherName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkWatcherName", *parsed)
	}

	if id.PacketCaptureName, ok = parsed.Parsed["packetCaptureName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "packetCaptureName", *parsed)
	}

	return &id, nil
}

// ValidatePacketCaptureID checks that 'input' can be parsed as a Packet Capture ID
func ValidatePacketCaptureID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePacketCaptureID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Packet Capture ID
func (id PacketCaptureId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/packetCaptures/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName, id.PacketCaptureName)
}

// Segments returns a slice of Resource ID Segments which comprise this Packet Capture ID
func (id PacketCaptureId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkWatchers", "networkWatchers", "networkWatchers"),
		resourceids.UserSpecifiedSegment("networkWatcherName", "networkWatcherValue"),
		resourceids.StaticSegment("staticPacketCaptures", "packetCaptures", "packetCaptures"),
		resourceids.UserSpecifiedSegment("packetCaptureName", "packetCaptureValue"),
	}
}

// String returns a human-readable description of this Packet Capture ID
func (id PacketCaptureId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Watcher Name: %q", id.NetworkWatcherName),
		fmt.Sprintf("Packet Capture Name: %q", id.PacketCaptureName),
	}
	return fmt.Sprintf("Packet Capture (%s)", strings.Join(components, "\n"))
}
