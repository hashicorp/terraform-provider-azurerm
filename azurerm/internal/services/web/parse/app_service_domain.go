package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceDomainId struct {
	ResourceGroup string
	Name          string
}

func AppServiceDomainID(input string) (*AppServiceDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing App Service Domain ID %q: %+v", input, err)
	}

	domainId := AppServiceDomainId{
		ResourceGroup: id.ResourceGroup,
	}

	if domainId.Name, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainId, nil
}
