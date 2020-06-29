package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KeyVaultId struct {
	Name          string
	ResourceGroup string
}

func KeyVaultID(input string) (*KeyVaultId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := KeyVaultId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
