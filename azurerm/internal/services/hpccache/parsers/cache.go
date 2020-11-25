package parsers

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CacheId struct {
	ResourceGroup string
	Name          string
}

func ParseCacheID(input string) (*CacheId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cache := CacheId{
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
