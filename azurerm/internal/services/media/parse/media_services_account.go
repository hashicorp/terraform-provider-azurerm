package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MediaServicesAccountId struct {
	ResourceGroup string
	Name          string
}

func MediaServicesAccountID(input string) (*MediaServicesAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Media Services Account ID %q: %+v", input, err)
	}

	service := MediaServicesAccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("mediaservices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
