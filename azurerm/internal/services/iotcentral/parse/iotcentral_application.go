package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IoTCentralApplicationId struct {
	ResourceGroup string
	Name          string
}

func IoTCentralApplicationID(input string) (*IoTCentralApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse IoT Central Application ID %q: %+v", input, err)
	}

	server := IoTCentralApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("IoTApps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
