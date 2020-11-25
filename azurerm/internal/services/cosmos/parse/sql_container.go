package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlContainerId struct {
	ResourceGroup       string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
}

func SqlContainerID(input string) (*SqlContainerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse SQL Container ID %q: %+v", input, err)
	}

	sqlContainer := SqlContainerId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlContainer.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if sqlContainer.SqlDatabaseName, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}

	if sqlContainer.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}

	return &sqlContainer, nil
}
