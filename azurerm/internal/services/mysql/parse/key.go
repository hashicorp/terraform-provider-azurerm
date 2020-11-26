package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KeyId struct {
	ResourceGroup string
	ServerName    string
	Name          string
}

func KeyID(input string) (*KeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse MySQL Server Key ID %q: %+v", input, err)
	}

	key := KeyId{
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
