package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualWANId{}

// VirtualWANId is a struct representing the Resource ID for a Virtual W A N
type VirtualWANId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualWanName    string
}

// NewVirtualWANID returns a new VirtualWANId struct
func NewVirtualWANID(subscriptionId string, resourceGroupName string, virtualWanName string) VirtualWANId {
	return VirtualWANId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualWanName:    virtualWanName,
	}
}

// ParseVirtualWANID parses 'input' into a VirtualWANId
func ParseVirtualWANID(input string) (*VirtualWANId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualWANId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualWANId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualWanName, ok = parsed.Parsed["virtualWanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualWanName", *parsed)
	}

	return &id, nil
}

// ParseVirtualWANIDInsensitively parses 'input' case-insensitively into a VirtualWANId
// note: this method should only be used for API response data and not user input
func ParseVirtualWANIDInsensitively(input string) (*VirtualWANId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualWANId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualWANId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualWanName, ok = parsed.Parsed["virtualWanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualWanName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualWANID checks that 'input' can be parsed as a Virtual W A N ID
func ValidateVirtualWANID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualWANID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual W A N ID
func (id VirtualWANId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualWans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualWanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual W A N ID
func (id VirtualWANId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualWans", "virtualWans", "virtualWans"),
		resourceids.UserSpecifiedSegment("virtualWanName", "virtualWanValue"),
	}
}

// String returns a human-readable description of this Virtual W A N ID
func (id VirtualWANId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Wan Name: %q", id.VirtualWanName),
	}
	return fmt.Sprintf("Virtual W A N (%s)", strings.Join(components, "\n"))
}
