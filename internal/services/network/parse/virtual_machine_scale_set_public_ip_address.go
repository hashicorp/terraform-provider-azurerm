// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualMachineScaleSetPublicIPAddressId struct {
	SubscriptionId             string
	ResourceGroup              string
	VirtualMachineScaleSetName string
	VirtualMachineName         string
	NetworkInterfaceName       string
	IpConfigurationName        string
	PublicIPAddressName        string
}

func NewVirtualMachineScaleSetPublicIPAddressID(subscriptionId, resourceGroup, virtualMachineScaleSetName, virtualMachineName, networkInterfaceName, ipConfigurationName, publicIPAddressName string) VirtualMachineScaleSetPublicIPAddressId {
	return VirtualMachineScaleSetPublicIPAddressId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		VirtualMachineName:         virtualMachineName,
		NetworkInterfaceName:       networkInterfaceName,
		IpConfigurationName:        ipConfigurationName,
		PublicIPAddressName:        publicIPAddressName,
	}
}

func (id VirtualMachineScaleSetPublicIPAddressId) String() string {
	segments := []string{
		fmt.Sprintf("Public I P Address Name %q", id.PublicIPAddressName),
		fmt.Sprintf("Ip Configuration Name %q", id.IpConfigurationName),
		fmt.Sprintf("Network Interface Name %q", id.NetworkInterfaceName),
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Virtual Machine Scale Set Name %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Machine Scale Set Public I P Address", segmentsStr)
}

func (id VirtualMachineScaleSetPublicIPAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/virtualMachines/%s/networkInterfaces/%s/ipConfigurations/%s/publicIPAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineScaleSetName, id.VirtualMachineName, id.NetworkInterfaceName, id.IpConfigurationName, id.PublicIPAddressName)
}

// VirtualMachineScaleSetPublicIPAddressID parses a VirtualMachineScaleSetPublicIPAddress ID into an VirtualMachineScaleSetPublicIPAddressId struct
func VirtualMachineScaleSetPublicIPAddressID(input string) (*VirtualMachineScaleSetPublicIPAddressId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualMachineScaleSetPublicIPAddress ID: %+v", input, err)
	}

	resourceId := VirtualMachineScaleSetPublicIPAddressId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualMachineScaleSetName, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}
	if resourceId.VirtualMachineName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}
	if resourceId.NetworkInterfaceName, err = id.PopSegment("networkInterfaces"); err != nil {
		return nil, err
	}
	if resourceId.IpConfigurationName, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}
	if resourceId.PublicIPAddressName, err = id.PopSegment("publicIPAddresses"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
