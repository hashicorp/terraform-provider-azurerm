package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppConfigurationId struct {
	ResourceGroup string
	Name          string
}

func AppConfigurationID(input string) (*AppConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Configuration Server ID %q: %+v", input, err)
	}

	server := AppConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("configurationStores"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
