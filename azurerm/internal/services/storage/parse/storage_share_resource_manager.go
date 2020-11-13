package parse

import (
	"fmt"
)

type StorageShareResourceManagerId struct {
	Name            string
	AccountName     string
	FileServiceName string
	ResourceGroup   string
}

func (id StorageShareResourceManagerId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/shares/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.AccountName, id.FileServiceName, id.Name)
}

func NewStorageShareResourceManagerId(resourceGroup, accountName, containerName string) StorageShareResourceManagerId {
	return StorageShareResourceManagerId{
		Name:            containerName,
		AccountName:     accountName,
		FileServiceName: "default",
		ResourceGroup:   resourceGroup,
	}
}
