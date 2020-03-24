package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CustomResourceProviderId struct {
	ResourceGroup string
	Name          string
}

func CustomResourceProviderID(input string) (*CustomResourceProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Custom Resource Provider ID %q: %+v", input, err)
	}

	service := CustomResourceProviderId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("resourceproviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
