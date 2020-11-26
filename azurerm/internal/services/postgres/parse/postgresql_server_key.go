package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PostgreSQLServerKeyId struct {
	Name          string
	ServerName    string
	ResourceGroup string
}

func PostgreSQLServerKeyID(input string) (*PostgreSQLServerKeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Postgres Server Key ID %q: %+v", input, err)
	}

	server := PostgreSQLServerKeyId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if server.Name, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
