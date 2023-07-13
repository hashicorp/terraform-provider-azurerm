// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualMachineScaleSetId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewVirtualMachineScaleSetID(subscriptionId, resourceGroup, name string) VirtualMachineScaleSetId {
	return VirtualMachineScaleSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id VirtualMachineScaleSetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Machine Scale Set", segmentsStr)
}

func (id VirtualMachineScaleSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// VirtualMachineScaleSetID parses a VirtualMachineScaleSet ID into an VirtualMachineScaleSetId struct
func VirtualMachineScaleSetID(input string) (*VirtualMachineScaleSetId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualMachineScaleSet ID: %+v", input, err)
	}

	resourceId := VirtualMachineScaleSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// VirtualMachineScaleSetIDInsensitively parses an VirtualMachineScaleSet ID into an VirtualMachineScaleSetId struct, insensitively
// This should only be used to parse an ID for rewriting, the VirtualMachineScaleSetID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VirtualMachineScaleSetIDInsensitively(input string) (*VirtualMachineScaleSetId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineScaleSetId{
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
	if resourceId.Name, err = id.PopSegment(virtualMachineScaleSetsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
