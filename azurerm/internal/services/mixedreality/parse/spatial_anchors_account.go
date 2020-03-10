package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpatialAnchorsAccountId struct {
	ResourceGroup string
	Name          string
}

func SpatialAnchorsAccountID(input string) (*SpatialAnchorsAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Search Service ID %q: %+v", input, err)
	}

	service := SpatialAnchorsAccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("spatialAnchorsAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
