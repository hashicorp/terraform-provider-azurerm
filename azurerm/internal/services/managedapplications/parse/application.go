package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	ResourceGroup string
	Name          string
}

func ApplicationID(input string) (*ApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := ApplicationId{
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
