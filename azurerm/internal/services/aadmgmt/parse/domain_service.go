package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainServiceId struct {
	ResourceGroup string
	Name          string
}

func DomainServiceID(input string) (*DomainServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Domain Service ID %q: %+v", input, err)
	}

	domainService := DomainServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if domainService.Name, err = id.PopSegment("domainServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainService, nil
}

