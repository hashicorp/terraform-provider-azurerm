package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PostgreSQLServerId struct {
	ResourceGroup string
	Name          string
}

func PostgreSQLServerID(input string) (*PostgreSQLServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PostgreSQL Server ID %q: %+v", input, err)
	}

	server := PostgreSQLServerId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
