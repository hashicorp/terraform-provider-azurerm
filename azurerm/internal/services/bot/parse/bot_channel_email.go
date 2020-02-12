package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BotChannelEmailId struct {
	ResourceGroup string
	BotName       string
}

func BotChannelEmailID(input string) (*BotChannelEmailId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Bot Channel Email ID %q: %+v", input, err)
	}

	service := BotChannelEmailId{
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
