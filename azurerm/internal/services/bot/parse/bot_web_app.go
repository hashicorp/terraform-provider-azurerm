package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BotWebAppId struct {
	ResourceGroup string
	Name          string
}

func BotWebAppID(input string) (*BotWebAppId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Bot Web App ID %q: %+v", input, err)
	}

	service := BotWebAppId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("botServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
