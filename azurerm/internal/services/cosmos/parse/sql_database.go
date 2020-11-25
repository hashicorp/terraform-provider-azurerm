package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlDatabaseId struct {
	ResourceGroup       string
	DatabaseAccountName string
	Name                string
}

func SqlDatabaseID(input string) (*SqlDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse SQL Database ID %q: %+v", input, err)
	}

	sqlDatabase := SqlDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlDatabase.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if sqlDatabase.Name, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}

	return &sqlDatabase, nil
}
