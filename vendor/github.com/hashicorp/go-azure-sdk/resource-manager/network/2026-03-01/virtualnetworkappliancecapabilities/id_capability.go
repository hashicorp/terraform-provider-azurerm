package virtualnetworkappliancecapabilities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CapabilityId{})
}

var _ resourceids.ResourceId = &CapabilityId{}

// CapabilityId is a struct representing the Resource ID for a Capability
type CapabilityId struct {
	SubscriptionId              string
	ResourceGroupName           string
	VirtualNetworkApplianceName string
	CapabilityName              string
}

// NewCapabilityID returns a new CapabilityId struct
func NewCapabilityID(subscriptionId string, resourceGroupName string, virtualNetworkApplianceName string, capabilityName string) CapabilityId {
	return CapabilityId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		VirtualNetworkApplianceName: virtualNetworkApplianceName,
		CapabilityName:              capabilityName,
	}
}

// ParseCapabilityID parses 'input' into a CapabilityId
func ParseCapabilityID(input string) (*CapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCapabilityIDInsensitively parses 'input' case-insensitively into a CapabilityId
// note: this method should only be used for API response data and not user input
func ParseCapabilityIDInsensitively(input string) (*CapabilityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CapabilityId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualNetworkApplianceName, ok = input.Parsed["virtualNetworkApplianceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkApplianceName", input)
	}

	if id.CapabilityName, ok = input.Parsed["capabilityName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capabilityName", input)
	}

	return nil
}

// ValidateCapabilityID checks that 'input' can be parsed as a Capability ID
func ValidateCapabilityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapabilityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capability ID
func (id CapabilityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkAppliances/%s/capabilities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkApplianceName, id.CapabilityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capability ID
func (id CapabilityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworkAppliances", "virtualNetworkAppliances", "virtualNetworkAppliances"),
		resourceids.UserSpecifiedSegment("virtualNetworkApplianceName", "virtualNetworkApplianceName"),
		resourceids.StaticSegment("staticCapabilities", "capabilities", "capabilities"),
		resourceids.UserSpecifiedSegment("capabilityName", "capabilityName"),
	}
}

// String returns a human-readable description of this Capability ID
func (id CapabilityId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Appliance Name: %q", id.VirtualNetworkApplianceName),
		fmt.Sprintf("Capability Name: %q", id.CapabilityName),
	}
	return fmt.Sprintf("Capability (%s)", strings.Join(components, "\n"))
}
