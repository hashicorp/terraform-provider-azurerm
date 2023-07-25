// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DataDiskId struct {
	SubscriptionId     string
	ResourceGroup      string
	VirtualMachineName string
	Name               string
}

func NewDataDiskID(subscriptionId, resourceGroup, virtualMachineName, name string) DataDiskId {
	return DataDiskId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		VirtualMachineName: virtualMachineName,
		Name:               name,
	}
}

func (id DataDiskId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Data Disk", segmentsStr)
}

func (id DataDiskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/dataDisks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName, id.Name)
}

// DataDiskID parses a DataDisk ID into an DataDiskId struct
func DataDiskID(input string) (*DataDiskId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DataDisk ID: %+v", input, err)
	}

	resourceId := DataDiskId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualMachineName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("dataDisks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
