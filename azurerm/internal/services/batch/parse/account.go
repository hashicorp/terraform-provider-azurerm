package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountId struct {
	ResourceGroup    string
	BatchAccountName string
}

func AccountID(input string) (*AccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Account ID %q: %+v", input, err)
	}

	account := AccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
