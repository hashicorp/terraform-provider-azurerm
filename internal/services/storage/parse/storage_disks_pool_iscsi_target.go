package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type StorageDisksPoolISCSITargetId struct {
	SubscriptionId  string
	ResourceGroup   string
	DiskPoolName    string
	IscsiTargetName string
}

func NewStorageDisksPoolISCSITargetID(subscriptionId, resourceGroup, diskPoolName, iscsiTargetName string) StorageDisksPoolISCSITargetId {
	return StorageDisksPoolISCSITargetId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		DiskPoolName:    diskPoolName,
		IscsiTargetName: iscsiTargetName,
	}
}

func (id StorageDisksPoolISCSITargetId) String() string {
	segments := []string{
		fmt.Sprintf("Iscsi Target Name %q", id.IscsiTargetName),
		fmt.Sprintf("Disk Pool Name %q", id.DiskPoolName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Disks Pool I S C S I Target", segmentsStr)
}

func (id StorageDisksPoolISCSITargetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StoragePool/diskPools/%s/iscsiTargets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
}

// StorageDisksPoolISCSITargetID parses a StorageDisksPoolISCSITarget ID into an StorageDisksPoolISCSITargetId struct
func StorageDisksPoolISCSITargetID(input string) (*StorageDisksPoolISCSITargetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageDisksPoolISCSITargetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DiskPoolName, err = id.PopSegment("diskPools"); err != nil {
		return nil, err
	}
	if resourceId.IscsiTargetName, err = id.PopSegment("iscsiTargets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
