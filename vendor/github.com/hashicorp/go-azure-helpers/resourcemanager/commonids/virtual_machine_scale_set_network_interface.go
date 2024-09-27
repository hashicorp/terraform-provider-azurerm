// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VirtualMachineScaleSetNetworkInterfaceId{}

// VirtualMachineScaleSetNetworkInterfaceId is a struct representing the Resource ID for a Virtual Machine Scale Set Network Interface
type VirtualMachineScaleSetNetworkInterfaceId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VirtualMachineScaleSetName string
	VirtualMachineIndex        string
	NetworkInterfaceName       string
}

// NewVirtualMachineScaleSetNetworkInterfaceID returns a new VirtualMachineScaleSetNetworkInterfaceId struct
func NewVirtualMachineScaleSetNetworkInterfaceID(subscriptionId string, resourceGroupName string, virtualMachineScaleSetName string, virtualMachineIndex string, networkInterfaceName string) VirtualMachineScaleSetNetworkInterfaceId {
	return VirtualMachineScaleSetNetworkInterfaceId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		VirtualMachineIndex:        virtualMachineIndex,
		NetworkInterfaceName:       networkInterfaceName,
	}
}

// ParseVirtualMachineScaleSetNetworkInterfaceID parses 'input' into a VirtualMachineScaleSetNetworkInterfaceId
func ParseVirtualMachineScaleSetNetworkInterfaceID(input string) (*VirtualMachineScaleSetNetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetNetworkInterfaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetNetworkInterfaceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineScaleSetNetworkInterfaceIDInsensitively parses 'input' case-insensitively into a VirtualMachineScaleSetNetworkInterfaceId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineScaleSetNetworkInterfaceIDInsensitively(input string) (*VirtualMachineScaleSetNetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetNetworkInterfaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetNetworkInterfaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineScaleSetNetworkInterfaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualMachineScaleSetName, ok = input.Parsed["virtualMachineScaleSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineScaleSetName", input)
	}

	if id.VirtualMachineIndex, ok = input.Parsed["virtualMachineIndex"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineIndex", input)
	}

	if id.NetworkInterfaceName, ok = input.Parsed["networkInterfaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", input)
	}

	return nil
}

// ValidateVirtualMachineScaleSetNetworkInterfaceID checks that 'input' can be parsed as a Virtual Machine Scale Set Network Interface ID
func ValidateVirtualMachineScaleSetNetworkInterfaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineScaleSetNetworkInterfaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Scale Set Network Interface ID
func (id VirtualMachineScaleSetNetworkInterfaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/virtualMachines/%s/networkInterfaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName, id.VirtualMachineIndex, id.NetworkInterfaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Scale Set Network Interface ID
func (id VirtualMachineScaleSetNetworkInterfaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("virtualMachineScaleSets", "virtualMachineScaleSets", "virtualMachineScaleSets"),
		resourceids.UserSpecifiedSegment("virtualMachineScaleSetName", "virtualMachineScaleSetValue"),
		resourceids.StaticSegment("virtualMachines", "virtualMachines", "virtualMachines"),
		resourceids.UserSpecifiedSegment("virtualMachineIndex", "virtualMachineIndexValue"),
		resourceids.StaticSegment("networkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
	}
}

// String returns a human-readable description of this Virtual Machine Scale Set Network Interface ID
func (id VirtualMachineScaleSetNetworkInterfaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Scale Set Name: %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Virtual Machine Index: %q", id.VirtualMachineIndex),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
	}
	return fmt.Sprintf("Virtual Machine Scale Set Network Interface (%s)", strings.Join(components, "\n"))
}
