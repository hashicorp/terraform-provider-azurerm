package virtualnetworkappliances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkApplianceId{})
}

var _ resourceids.ResourceId = &VirtualNetworkApplianceId{}

// VirtualNetworkApplianceId is a struct representing the Resource ID for a Virtual Network Appliance
type VirtualNetworkApplianceId struct {
	SubscriptionId              string
	ResourceGroupName           string
	VirtualNetworkApplianceName string
}

// NewVirtualNetworkApplianceID returns a new VirtualNetworkApplianceId struct
func NewVirtualNetworkApplianceID(subscriptionId string, resourceGroupName string, virtualNetworkApplianceName string) VirtualNetworkApplianceId {
	return VirtualNetworkApplianceId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		VirtualNetworkApplianceName: virtualNetworkApplianceName,
	}
}

// ParseVirtualNetworkApplianceID parses 'input' into a VirtualNetworkApplianceId
func ParseVirtualNetworkApplianceID(input string) (*VirtualNetworkApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkApplianceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkApplianceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkApplianceIDInsensitively parses 'input' case-insensitively into a VirtualNetworkApplianceId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkApplianceIDInsensitively(input string) (*VirtualNetworkApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkApplianceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkApplianceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkApplianceId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateVirtualNetworkApplianceID checks that 'input' can be parsed as a Virtual Network Appliance ID
func ValidateVirtualNetworkApplianceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkApplianceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Appliance ID
func (id VirtualNetworkApplianceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkAppliances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkApplianceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Appliance ID
func (id VirtualNetworkApplianceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworkAppliances", "virtualNetworkAppliances", "virtualNetworkAppliances"),
		resourceids.UserSpecifiedSegment("virtualNetworkApplianceName", "virtualNetworkApplianceName"),
	}
}

// String returns a human-readable description of this Virtual Network Appliance ID
func (id VirtualNetworkApplianceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Appliance Name: %q", id.VirtualNetworkApplianceName),
	}
	return fmt.Sprintf("Virtual Network Appliance (%s)", strings.Join(components, "\n"))
}
