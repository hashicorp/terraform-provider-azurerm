package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SearchServiceId struct {
	ResourceGroup string
	Name          string
}

func SearchServiceID(input string) (*SearchServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Search Service ID %q: %+v", input, err)
	}

	service := SearchServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("searchServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
