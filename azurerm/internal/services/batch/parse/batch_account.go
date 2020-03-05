package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BatchAccountId struct {
	ResourceGroup string
	Name          string
}

func BatchAccountID(input string) (*BatchAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Account ID %q: %+v", input, err)
	}

	account := BatchAccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
