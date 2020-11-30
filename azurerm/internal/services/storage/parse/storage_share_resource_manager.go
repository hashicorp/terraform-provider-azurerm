package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageShareResourceManagerId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	FileServiceName    string
	ShareName          string
}

func NewStorageShareResourceManagerID(subscriptionId, resourceGroup, storageAccountName, fileServiceName, shareName string) StorageShareResourceManagerId {
	return StorageShareResourceManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		FileServiceName:    fileServiceName,
		ShareName:          shareName,
	}
}

func (id StorageShareResourceManagerId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/shares/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileServiceName, id.ShareName)
}

// StorageShareResourceManagerID parses a StorageShareResourceManager ID into an StorageShareResourceManagerId struct
func StorageShareResourceManagerID(input string) (*StorageShareResourceManagerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageShareResourceManagerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.FileServiceName, err = id.PopSegment("fileServices"); err != nil {
		return nil, err
	}
	if resourceId.ShareName, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
