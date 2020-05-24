package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedApplicationDefinitionId struct {
	Name          string
	ResourceGroup string
}

func ManagedApplicationDefinitionID(input string) (*ManagedApplicationDefinitionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := ManagedApplicationDefinitionId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("applicationDefinitions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
