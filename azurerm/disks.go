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

	blobDomainSuffix := meta.(*ArmClient).environment.StorageEndpointSuffix
	if !strings.HasSuffix(strings.ToLower(vhdURL.Host), strings.ToLower(blobDomainSuffix)) {
		return nil, fmt.Errorf("Error: Disk VHD URI %q doesn't appear to be a Blob Storage URI (%q) - expected a suffix of %q)", uri, vhdURL.Host, blobDomainSuffix)
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
