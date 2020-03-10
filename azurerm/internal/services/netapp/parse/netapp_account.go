package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetAppAccountId struct {
	ResourceGroup string
	Name          string
}

func NetAppAccountID(input string) (*NetAppAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Account ID %q: %+v", input, err)
	}

	service := NetAppAccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
