// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VMSSInstanceId struct {
	SubscriptionId             string
	ResourceGroup              string
	VirtualMachineScaleSetName string
	VirtualMachineName         string
}

func NewVMSSInstanceID(subscriptionId, resourceGroup, virtualMachineScaleSetName, virtualMachineName string) VMSSInstanceId {
	return VMSSInstanceId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		VirtualMachineName:         virtualMachineName,
	}
}

func (id VMSSInstanceId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Virtual Machine Scale Set Name %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "V M S S Instance", segmentsStr)
}

func (id VMSSInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/virtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineScaleSetName, id.VirtualMachineName)
}

// VMSSInstanceID parses a VMSSInstance ID into an VMSSInstanceId struct
func VMSSInstanceID(input string) (*VMSSInstanceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VMSSInstance ID: %+v", input, err)
	}

	resourceId := VMSSInstanceId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// VMSSInstanceIDInsensitively parses an VMSSInstance ID into an VMSSInstanceId struct, insensitively
// This should only be used to parse an ID for rewriting, the VMSSInstanceID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VMSSInstanceIDInsensitively(input string) (*VMSSInstanceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VMSSInstanceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'virtualMachineScaleSets' segment
	virtualMachineScaleSetsKey := "virtualMachineScaleSets"
	for key := range id.Path {
		if strings.EqualFold(key, virtualMachineScaleSetsKey) {
			virtualMachineScaleSetsKey = key
			break
		}
	}
	if resourceId.VirtualMachineScaleSetName, err = id.PopSegment(virtualMachineScaleSetsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'virtualMachines' segment
	virtualMachinesKey := "virtualMachines"
	for key := range id.Path {
		if strings.EqualFold(key, virtualMachinesKey) {
			virtualMachinesKey = key
			break
		}
	}
	if resourceId.VirtualMachineName, err = id.PopSegment(virtualMachinesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
