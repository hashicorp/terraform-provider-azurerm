package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	ResourceGroup string
	IoTAppName    string
}

func ApplicationID(input string) (*ApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse IoT Central Application ID %q: %+v", input, err)
	}

	server := ApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.IoTAppName, err = id.PopSegment("IoTApps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
