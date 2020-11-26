package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MediaServiceId struct {
	ResourceGroup string
	Name          string
}

func MediaServiceID(input string) (*MediaServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Media Services Account ID %q: %+v", input, err)
	}

	service := MediaServiceId{
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
