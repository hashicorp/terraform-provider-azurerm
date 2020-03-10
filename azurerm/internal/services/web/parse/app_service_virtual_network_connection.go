package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceVirtualNetworkConnectionId struct {
	ResourceGroup  string
	AppServiceName string
	Name           string
}

func AppServiceVirtualNetworkConnectionID(input string) (*AppServiceVirtualNetworkConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Environment ID %q: %+v", input, err)
	}

	appServiceVirtualNetworkConnection := AppServiceVirtualNetworkConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if appServiceVirtualNetworkConnection.AppServiceName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if appServiceVirtualNetworkConnection.Name, err = id.PopSegment("virtualNetworkConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceVirtualNetworkConnection, nil
}
