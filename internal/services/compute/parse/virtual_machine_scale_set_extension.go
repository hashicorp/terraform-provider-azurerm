// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualMachineScaleSetExtensionId struct {
	SubscriptionId             string
	ResourceGroup              string
	VirtualMachineScaleSetName string
	ExtensionName              string
}

func NewVirtualMachineScaleSetExtensionID(subscriptionId, resourceGroup, virtualMachineScaleSetName, extensionName string) VirtualMachineScaleSetExtensionId {
	return VirtualMachineScaleSetExtensionId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		ExtensionName:              extensionName,
	}
}

func (id VirtualMachineScaleSetExtensionId) String() string {
	segments := []string{
		fmt.Sprintf("Extension Name %q", id.ExtensionName),
		fmt.Sprintf("Virtual Machine Scale Set Name %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Machine Scale Set Extension", segmentsStr)
}

func (id VirtualMachineScaleSetExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName)
}

// VirtualMachineScaleSetExtensionID parses a VirtualMachineScaleSetExtension ID into an VirtualMachineScaleSetExtensionId struct
func VirtualMachineScaleSetExtensionID(input string) (*VirtualMachineScaleSetExtensionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualMachineScaleSetExtension ID: %+v", input, err)
	}

	resourceId := VirtualMachineScaleSetExtensionId{
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
	if resourceId.ExtensionName, err = id.PopSegment("extensions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
