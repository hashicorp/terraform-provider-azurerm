package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountId struct {
	ResourceGroup     string
	NetAppAccountName string
}

func NetAppAccountID(input string) (*AccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Account ID %q: %+v", input, err)
	}

	service := AccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
