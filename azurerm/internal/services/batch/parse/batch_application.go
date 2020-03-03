package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BatchApplicationId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func BatchApplicationID(input string) (*BatchApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Application ID %q: %+v", input, err)
	}

	application := BatchApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if application.AccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if application.Name, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &application, nil
}
