package parsers

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type StorageContainerResourceManagerId struct {
	Name            string
	AccountName     string
	BlobServiceName string
	ResourceGroup   string
}

func StorageContainerResourceManagerID(input string) (*StorageContainerResourceManagerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cache := StorageContainerResourceManagerId{
		ResourceGroup: id.ResourceGroup,
	}

	if cache.Name, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}

	if cache.BlobServiceName, err = id.PopSegment("blobServices"); err != nil {
		return nil, err
	}

	if cache.AccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cache, nil
}
