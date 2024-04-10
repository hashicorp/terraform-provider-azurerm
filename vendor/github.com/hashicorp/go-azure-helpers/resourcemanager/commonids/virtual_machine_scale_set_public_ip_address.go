// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VirtualMachineScaleSetPublicIPAddressId{}

// VirtualMachineScaleSetPublicIPAddressId is a struct representing the Resource ID for a Virtual Machine Scale Set Public I P Address
type VirtualMachineScaleSetPublicIPAddressId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VirtualMachineScaleSetName string
	VirtualMachineIndex        string
	NetworkInterfaceName       string
	IpConfigurationName        string
	PublicIpAddressName        string
}

// NewVirtualMachineScaleSetPublicIPAddressID returns a new VirtualMachineScaleSetPublicIPAddressId struct
func NewVirtualMachineScaleSetPublicIPAddressID(subscriptionId string, resourceGroupName string, virtualMachineScaleSetName string, virtualMachineIndex string, networkInterfaceName string, ipConfigurationName string, publicIpAddressName string) VirtualMachineScaleSetPublicIPAddressId {
	return VirtualMachineScaleSetPublicIPAddressId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		VirtualMachineIndex:        virtualMachineIndex,
		NetworkInterfaceName:       networkInterfaceName,
		IpConfigurationName:        ipConfigurationName,
		PublicIpAddressName:        publicIpAddressName,
	}
}

// ParseVirtualMachineScaleSetPublicIPAddressID parses 'input' into a VirtualMachineScaleSetPublicIPAddressId
func ParseVirtualMachineScaleSetPublicIPAddressID(input string) (*VirtualMachineScaleSetPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetPublicIPAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetPublicIPAddressId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineScaleSetPublicIPAddressIDInsensitively parses 'input' case-insensitively into a VirtualMachineScaleSetPublicIPAddressId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineScaleSetPublicIPAddressIDInsensitively(input string) (*VirtualMachineScaleSetPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetPublicIPAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetPublicIPAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineScaleSetPublicIPAddressId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PublicIpAddressName, ok = input.Parsed["publicIpAddressName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publicIpAddressName", input)
	}

	return nil
}

// ValidateVirtualMachineScaleSetPublicIPAddressID checks that 'input' can be parsed as a Virtual Machine Scale Set Public I P Address ID
func ValidateVirtualMachineScaleSetPublicIPAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineScaleSetPublicIPAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Scale Set Public IP Address ID
func (id VirtualMachineScaleSetPublicIPAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/virtualMachines/%s/networkInterfaces/%s/ipConfigurations/%s/publicIPAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName, id.VirtualMachineIndex, id.NetworkInterfaceName, id.IpConfigurationName, id.PublicIpAddressName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Scale Set Public I P Address ID
func (id VirtualMachineScaleSetPublicIPAddressId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("publicIPAddresses", "publicIPAddresses", "publicIPAddresses"),
		resourceids.UserSpecifiedSegment("publicIpAddressName", "publicIpAddressValue"),
	}
}

// String returns a human-readable description of this Virtual Machine Scale Set Public I P Address ID
func (id VirtualMachineScaleSetPublicIPAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Scale Set Name: %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Virtual Machine Index: %q", id.VirtualMachineIndex),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Ip Configuration Name: %q", id.IpConfigurationName),
		fmt.Sprintf("Public Ip Address Name: %q", id.PublicIpAddressName),
	}
	return fmt.Sprintf("Virtual Machine Scale Set Public IP Address (%s)", strings.Join(components, "\n"))
}
