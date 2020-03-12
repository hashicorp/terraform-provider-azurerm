package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataBoxJobId struct {
	Name          string
	ResourceGroup string
}

func DataBoxJobID(input string) (*DataBoxJobId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := DataBoxJobId{
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
