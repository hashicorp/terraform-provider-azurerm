package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseAccountId struct {
	ResourceGroup string
	Name          string
}

func DatabaseAccountID(input string) (*DatabaseAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Account ID %q: %+v", input, err)
	}

	databaseAccount := DatabaseAccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if databaseAccount.Name, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	return &databaseAccount, nil
}
