// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualRouterPeeringId{}

// VirtualRouterPeeringId is a struct representing the Resource ID for a Virtual Router Peering
type VirtualRouterPeeringId struct {
	SubscriptionId    string
	ResourceGroupName string
	VirtualRouterName string
	PeeringName       string
}

// NewVirtualRouterPeeringID returns a new VirtualRouterPeeringId struct
func NewVirtualRouterPeeringID(subscriptionId string, resourceGroupName string, virtualRouterName string, peeringName string) VirtualRouterPeeringId {
	return VirtualRouterPeeringId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VirtualRouterName: virtualRouterName,
		PeeringName:       peeringName,
	}
}

// ParseVirtualRouterPeeringID parses 'input' into a VirtualRouterPeeringId
func ParseVirtualRouterPeeringID(input string) (*VirtualRouterPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualRouterPeeringId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualRouterPeeringId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualRouterName, ok = parsed.Parsed["virtualRouterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualRouterName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ParseVirtualRouterPeeringIDInsensitively parses 'input' case-insensitively into a VirtualRouterPeeringId
// note: this method should only be used for API response data and not user input
func ParseVirtualRouterPeeringIDInsensitively(input string) (*VirtualRouterPeeringId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualRouterPeeringId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualRouterPeeringId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualRouterName, ok = parsed.Parsed["virtualRouterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualRouterName", *parsed)
	}

	if id.PeeringName, ok = parsed.Parsed["peeringName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "peeringName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualRouterPeeringID checks that 'input' can be parsed as a Virtual Router Peering ID
func ValidateVirtualRouterPeeringID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualRouterPeeringID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Router Peering ID
func (id VirtualRouterPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualRouters/%s/peerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualRouterName, id.PeeringName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Router Peering ID
func (id VirtualRouterPeeringId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("virtualRouters", "virtualRouters", "virtualRouters"),
		resourceids.UserSpecifiedSegment("virtualRouterName", "virtualRouterValue"),
		resourceids.StaticSegment("peerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringValue"),
	}
}

// String returns a human-readable description of this Virtual Router Peering ID
func (id VirtualRouterPeeringId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Router Name: %q", id.VirtualRouterName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
	}
	return fmt.Sprintf("Virtual Router Peering (%s)", strings.Join(components, "\n"))
}
