package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventGridDomainId struct {
	ResourceGroup string
	Name          string
}

func EventGridDomainID(input string) (*EventGridDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Domain ID %q: %+v", input, err)
	}

	domain := EventGridDomainId{
		ResourceGroup: id.ResourceGroup,
	}

	if domain.Name, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domain, nil
}
