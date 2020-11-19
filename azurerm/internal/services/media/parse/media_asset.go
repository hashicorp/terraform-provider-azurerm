package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MediaAssetstId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func MediaAssetsID(input string) (*MediaAssetstId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Asset ID %q: %+v", input, err)
	}

	asset := MediaAssetstId{
		ResourceGroup: id.ResourceGroup,
	}

	if asset.AccountName, err = id.PopSegment("mediaservices"); err != nil {
		return nil, err
	}

	if asset.Name, err = id.PopSegment("assets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &asset, nil
}
