package parse

import (
	"fmt"
)

type StorageShareResourceManagerId struct {
	ResourceGroup      string
	StorageAccountName string
	FileServiceName    string
	ShareName          string
}

func NewStorageShareResourceManagerID(resourceGroup, storageAccountName, fileServiceName, shareName string) StorageShareResourceManagerId {
	return StorageShareResourceManagerId{
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		FileServiceName:    fileServiceName,
		ShareName:          shareName,
	}
}

func (id StorageShareResourceManagerId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/shares/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileServiceName, id.ShareName)
}
