package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BotConnectionId struct {
	ResourceGroup string
	BotName       string
	Name          string
}

func BotConnectionID(input string) (*BotConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Bot Connection ID %q: %+v", input, err)
	}

	connection := BotConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if connection.BotName, err = id.PopSegment("botServices"); err != nil {
		return nil, err
	}

	if connection.Name, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &connection, nil
}
