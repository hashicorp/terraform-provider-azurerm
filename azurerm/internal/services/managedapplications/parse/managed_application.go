package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedApplicationId struct {
	Name          string
	ResourceGroup string
}

func ManagedApplicationID(input string) (*ManagedApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := ManagedApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
