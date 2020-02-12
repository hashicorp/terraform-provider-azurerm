package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BotChannelSlackId struct {
	ResourceGroup string
	BotName       string
	Name          string
}

func BotChannelSlackID(input string) (*BotChannelSlackId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Bot Channel Slack ID %q: %+v", input, err)
	}

	service := BotChannelSlackId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.BotName, err = id.PopSegment("botServices"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("channels"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
