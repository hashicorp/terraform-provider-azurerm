package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageSyncServiceId struct {
	ResourceGroup string
	Name          string
}

func (id StorageSyncServiceId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func StorageSyncServiceID(input string) (*StorageSyncServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	storageSync := StorageSyncServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if storageSync.Name, err = id.PopSegment("storageSyncServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &storageSync, nil
}
