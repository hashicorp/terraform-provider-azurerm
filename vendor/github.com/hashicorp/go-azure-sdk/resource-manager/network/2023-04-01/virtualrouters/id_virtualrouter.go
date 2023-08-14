package virtualrouters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualRouterId{}

// VirtualRouterId is a struct representing the Resource ID for a Virtual Router
type VirtualRouterId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualRouterName string
}

// NewVirtualRouterID returns a new VirtualRouterId struct
func NewVirtualRouterID(subscriptionId string, resourceGroupName string, virtualRouterName string) VirtualRouterId {
	return VirtualRouterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualRouterName: virtualRouterName,
	}
}

// ParseVirtualRouterID parses 'input' into a VirtualRouterId
func ParseVirtualRouterID(input string) (*VirtualRouterId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualRouterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualRouterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualRouterName, ok = parsed.Parsed["virtualRouterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualRouterName", *parsed)
	}

	return &id, nil
}

// ParseVirtualRouterIDInsensitively parses 'input' case-insensitively into a VirtualRouterId
// note: this method should only be used for API response data and not user input
func ParseVirtualRouterIDInsensitively(input string) (*VirtualRouterId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualRouterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualRouterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualRouterName, ok = parsed.Parsed["virtualRouterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualRouterName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualRouterID checks that 'input' can be parsed as a Virtual Router ID
func ValidateVirtualRouterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualRouterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Router ID
func (id VirtualRouterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualRouters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualRouterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Router ID
func (id VirtualRouterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualRouters", "virtualRouters", "virtualRouters"),
		resourceids.UserSpecifiedSegment("virtualRouterName", "virtualRouterValue"),
	}
}

// String returns a human-readable description of this Virtual Router ID
func (id VirtualRouterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Router Name: %q", id.VirtualRouterName),
	}
	return fmt.Sprintf("Virtual Router (%s)", strings.Join(components, "\n"))
}
