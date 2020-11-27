package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServerKeyId struct {
	ResourceGroup string
	ServerName    string
	KeyName       string
}

func ServerKeyID(input string) (*ServerKeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Postgres Server Key ID %q: %+v", input, err)
	}

	server := ServerKeyId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if server.KeyName, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
