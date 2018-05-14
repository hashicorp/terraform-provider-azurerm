package azurerm

import (
	"fmt"
	"net/url"
	"strings"
)

type UnmanagedDiskMetadata struct {
	StorageContainerName string
	StorageAccountName   string
	BlobName             string
	ResourceGroupName    string
}

func storageAccountNameForUnmanagedDisk(uri string, meta interface{}) (*UnmanagedDiskMetadata, error) {
	vhdURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Disk VHD URI: %s", err)
	}

	// VHD URI is in the form: https://storageAccountName.blob.core.windows.net/containerName/blobName
	storageAccountName := strings.Split(vhdURL.Host, ".")[0]
	path := strings.Split(strings.TrimPrefix(vhdURL.Path, "/"), "/")
	containerName := path[0]
	blobName := path[1]

	resourceGroupName, err := findStorageAccountResourceGroup(meta, storageAccountName)
	if err != nil {
		return nil, fmt.Errorf("Error finding resource group for storage account %s: %+v", storageAccountName, err)
	}

	metadata := UnmanagedDiskMetadata{
		BlobName:             blobName,
		StorageContainerName: containerName,
		StorageAccountName:   storageAccountName,
		ResourceGroupName:    resourceGroupName,
	}
	return &metadata, nil
}
