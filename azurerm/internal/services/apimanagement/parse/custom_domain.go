package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CustomDomainId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func CustomDomainID(input string) (*CustomDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	logger := CustomDomainId{
		ResourceGroup: id.ResourceGroup,
	}

	if logger.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if logger.Name, err = id.PopSegment("customDomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logger, nil
}
