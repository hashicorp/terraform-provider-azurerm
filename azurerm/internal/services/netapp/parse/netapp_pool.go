package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetAppPoolId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func NetAppPoolID(input string) (*NetAppPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Pool ID %q: %+v", input, err)
	}

	service := NetAppPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.AccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
