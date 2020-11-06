package parsers

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type StorageSyncGroupId struct {
	Name            string
	StorageSyncName string
	ResourceGroup   string
}

func StorageSyncGroupID(input string) (*StorageSyncGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	storageSyncGroup := StorageSyncGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if storageSyncGroup.StorageSyncName, err = id.PopSegment("storageSyncServices"); err != nil {
		return nil, err
	}

	if storageSyncGroup.Name, err = id.PopSegment("syncGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &storageSyncGroup, nil
}
