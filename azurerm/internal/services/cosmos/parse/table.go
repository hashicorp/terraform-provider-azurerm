package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TableId struct {
	ResourceGroup       string
	DatabaseAccountName string
	Name                string
}

func TableID(input string) (*TableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Table ID %q: %+v", input, err)
	}

	table := TableId{
		ResourceGroup: id.ResourceGroup,
	}

	if table.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if table.Name, err = id.PopSegment("tables"); err != nil {
		return nil, err
	}

	return &table, nil
}
