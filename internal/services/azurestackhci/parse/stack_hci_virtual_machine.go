// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StackHCIVirtualMachineId struct {
	SubscriptionId             string
	ResourceGroup              string
	MachineName                string
	VirtualMachineInstanceName string
}

func NewStackHCIVirtualMachineID(subscriptionId, resourceGroup, machineName, virtualMachineInstanceName string) StackHCIVirtualMachineId {
	return StackHCIVirtualMachineId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		MachineName:                machineName,
		VirtualMachineInstanceName: virtualMachineInstanceName,
	}
}

func (id StackHCIVirtualMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Machine Instance Name %q", id.VirtualMachineInstanceName),
		fmt.Sprintf("Machine Name %q", id.MachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "StackHCI Virtual Machine", segmentsStr)
}

func (id StackHCIVirtualMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/providers/Microsoft.AzureStackHCI/virtualMachineInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MachineName, id.VirtualMachineInstanceName)
}

// StackHCIVirtualMachineID parses a StackHCIVirtualMachine ID into an StackHCIVirtualMachineId struct
func StackHCIVirtualMachineID(input string) (*StackHCIVirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StackHCIVirtualMachine ID: %+v", input, err)
	}

	resourceId := StackHCIVirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MachineName, err = id.PopSegment("machines"); err != nil {
		return nil, err
	}
	if resourceId.VirtualMachineInstanceName, err = id.PopSegment("virtualMachineInstances"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// StackHCIVirtualMachineIDInsensitively parses an StackHCIVirtualMachine ID into an StackHCIVirtualMachineId struct, insensitively
// This should only be used to parse an ID for rewriting, the StackHCIVirtualMachineID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func StackHCIVirtualMachineIDInsensitively(input string) (*StackHCIVirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StackHCIVirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'machines' segment
	machinesKey := "machines"
	for key := range id.Path {
		if strings.EqualFold(key, machinesKey) {
			machinesKey = key
			break
		}
	}
	if resourceId.MachineName, err = id.PopSegment(machinesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'virtualMachineInstances' segment
	virtualMachineInstancesKey := "virtualMachineInstances"
	for key := range id.Path {
		if strings.EqualFold(key, virtualMachineInstancesKey) {
			virtualMachineInstancesKey = key
			break
		}
	}
	if resourceId.VirtualMachineInstanceName, err = id.PopSegment(virtualMachineInstancesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
