package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedDiskId struct {
	SubscriptionId string
	ResourceGroup  string
	DiskName       string
}

func NewManagedDiskID(subscriptionId, resourceGroup, diskName string) ManagedDiskId {
	return ManagedDiskId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DiskName:       diskName,
	}
}

func (id ManagedDiskId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DiskName)
}

// ManagedDiskID parses a ManagedDisk ID into an ManagedDiskId struct
func ManagedDiskID(input string) (*ManagedDiskId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedDiskId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DiskName, err = id.PopSegment("disks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
