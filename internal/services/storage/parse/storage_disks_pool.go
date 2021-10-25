package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type StorageDisksPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	DiskPoolName   string
}

func NewStorageDisksPoolID(subscriptionId, resourceGroup, diskPoolName string) StorageDisksPoolId {
	return StorageDisksPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DiskPoolName:   diskPoolName,
	}
}

func (id StorageDisksPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Disk Pool Name %q", id.DiskPoolName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Disks Pool", segmentsStr)
}

func (id StorageDisksPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StoragePool/diskPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
}

// StorageDisksPoolID parses a StorageDisksPool ID into an StorageDisksPoolId struct
func StorageDisksPoolID(input string) (*StorageDisksPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageDisksPoolId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
