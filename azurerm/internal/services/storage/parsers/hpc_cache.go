package parsers

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HPCCacheID struct {
	Name          string
	ResourceGroup string
}

func ParseHPCCacheID(input string) (*HPCCacheID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cache := HPCCacheID{
		ResourceGroup: id.ResourceGroup,
	}

	if cache.Name, err = id.PopSegment("caches"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cache, nil
}
