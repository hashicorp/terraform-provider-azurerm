package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualHubId{}

// VirtualHubId is a struct representing the Resource ID for a Virtual Hub
type VirtualHubId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualHubName    string
}

// NewVirtualHubID returns a new VirtualHubId struct
func NewVirtualHubID(subscriptionId string, resourceGroupName string, virtualHubName string) VirtualHubId {
	return VirtualHubId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualHubName:    virtualHubName,
	}
}

// ParseVirtualHubID parses 'input' into a VirtualHubId
func ParseVirtualHubID(input string) (*VirtualHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualHubName, ok = parsed.Parsed["virtualHubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", *parsed)
	}

	return &id, nil
}

// ParseVirtualHubIDInsensitively parses 'input' case-insensitively into a VirtualHubId
// note: this method should only be used for API response data and not user input
func ParseVirtualHubIDInsensitively(input string) (*VirtualHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualHubName, ok = parsed.Parsed["virtualHubName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualHubID checks that 'input' can be parsed as a Virtual Hub ID
func ValidateVirtualHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Hub ID
func (id VirtualHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Hub ID
func (id VirtualHubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubValue"),
	}
}

// String returns a human-readable description of this Virtual Hub ID
func (id VirtualHubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
	}
	return fmt.Sprintf("Virtual Hub (%s)", strings.Join(components, "\n"))
}
