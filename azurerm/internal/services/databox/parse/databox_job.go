package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataBoxJobID struct {
	Name          string
	ResourceGroup string
}

func ParseDataBoxJobID(input string) (*DataBoxJobID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := DataBoxJobID{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("jobs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
