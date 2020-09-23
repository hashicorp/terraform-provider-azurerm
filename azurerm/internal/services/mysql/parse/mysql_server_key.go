package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MySQLServerKeyId struct {
	ResourceGroup string
	ServerName    string
	Name          string
}

func MySQLServerKeyID(input string) (*MySQLServerKeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse MySQL Server Key ID %q: %+v", input, err)
	}

	key := MySQLServerKeyId{
		ResourceGroup: id.ResourceGroup,
	}

	if key.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if key.Name, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &key, nil
}
