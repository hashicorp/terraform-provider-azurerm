package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VaultId struct {
	ResourceGroup string
	Name          string
}

func VaultID(input string) (*VaultId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := VaultId{
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
