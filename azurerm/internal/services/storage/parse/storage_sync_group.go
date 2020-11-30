package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type StorageSyncGroupId struct {
	ResourceGroup          string
	StorageSyncServiceName string
	SyncGroupName          string
}

func StorageSyncGroupID(input string) (*StorageSyncGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	storageSyncGroup := StorageSyncGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if storageSyncGroup.StorageSyncServiceName, err = id.PopSegment("storageSyncServices"); err != nil {
		return nil, err
	}

	if storageSyncGroup.SyncGroupName, err = id.PopSegment("syncGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &storageSyncGroup, nil
}
