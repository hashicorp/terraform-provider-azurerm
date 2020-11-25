package parsers

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageTargetId struct {
	ResourceGroup string
	Cache         string
	Name          string
}

func StorageTargetID(input string) (*StorageTargetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse HPC Cache Target ID %q: %+v", input, err)
	}

	target := StorageTargetId{
		ResourceGroup: id.ResourceGroup,
	}

	if target.Cache, err = id.PopSegment("caches"); err != nil {
		return nil, err
	}

	if target.Name, err = id.PopSegment("storageTargets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &target, nil
}
