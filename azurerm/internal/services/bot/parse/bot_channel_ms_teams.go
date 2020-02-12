package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BotChannelMsTeamsId struct {
	ResourceGroup string
	BotName       string
}

func BotChannelMsTeamsID(input string) (*BotChannelMsTeamsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Bot Channel Ms Teams ID %q: %+v", input, err)
	}

	service := BotChannelMsTeamsId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.BotName, err = id.PopSegment("botServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
