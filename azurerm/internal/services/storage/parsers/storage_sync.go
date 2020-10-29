package parsers

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageSyncId struct {
	Name          string
	ResourceGroup string
}

func (id StorageSyncId) ID(subscriptionId string) string {
	fmtString := ""
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func ParseStorageSyncID(input string) (*StorageSyncId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	storageSync := StorageSyncId{
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
