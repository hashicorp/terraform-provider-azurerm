// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VirtualMachineScaleSetIPConfigurationId{}

// VirtualMachineScaleSetIPConfigurationId is a struct representing the Resource ID for a Virtual Machine Scale Set Public I P Address
type VirtualMachineScaleSetIPConfigurationId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VirtualMachineScaleSetName string
	VirtualMachineIndex        string
	NetworkInterfaceName       string
	IpConfigurationName        string
}

// NewVirtualMachineScaleSetIPConfigurationId returns a new VirtualMachineScaleSetIPConfigurationId struct
func NewVirtualMachineScaleSetIPConfigurationID(subscriptionId string, resourceGroupName string, virtualMachineScaleSetName string, virtualMachineIndex string, networkInterfaceName string, ipConfigurationName string) VirtualMachineScaleSetIPConfigurationId {
	return VirtualMachineScaleSetIPConfigurationId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		VirtualMachineIndex:        virtualMachineIndex,
		NetworkInterfaceName:       networkInterfaceName,
		IpConfigurationName:        ipConfigurationName,
	}
}

// ParseVirtualMachineScaleSetIPConfigurationId parses 'input' into a VirtualMachineScaleSetIPConfigurationId
func ParseVirtualMachineScaleSetIPConfigurationId(input string) (*VirtualMachineScaleSetIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetIPConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetIPConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineScaleSetIPConfigurationIdInsensitively parses 'input' case-insensitively into a VirtualMachineScaleSetIPConfigurationId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineScaleSetIPConfigurationIdInsensitively(input string) (*VirtualMachineScaleSetIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetIPConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetIPConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineScaleSetIPConfigurationId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.IpConfigurationName, ok = input.Parsed["ipConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ipConfigurationName", input)
	}

	return nil
}

// ValidateVirtualMachineScaleSetIPConfigurationId checks that 'input' can be parsed as a Virtual Machine Scale Set Public I P Address ID
func ValidateVirtualMachineScaleSetIPConfigurationId(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineScaleSetIPConfigurationId(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Scale Set Public I P Address ID
func (id VirtualMachineScaleSetIPConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/virtualMachines/%s/networkInterfaces/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName, id.VirtualMachineIndex, id.NetworkInterfaceName, id.IpConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Scale Set Public I P Address ID
func (id VirtualMachineScaleSetIPConfigurationId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("ipConfigurations", "ipConfigurations", "ipConfigurations"),
		resourceids.UserSpecifiedSegment("ipConfigurationName", "ipConfigurationValue"),
	}
}

// String returns a human-readable description of this Virtual Machine Scale Set Public I P Address ID
func (id VirtualMachineScaleSetIPConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Scale Set Name: %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Virtual Machine Index: %q", id.VirtualMachineIndex),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Ip Configuration Name: %q", id.IpConfigurationName),
	}
	return fmt.Sprintf("Virtual Machine Scale Set IP Configuration (%s)", strings.Join(components, "\n"))
}
