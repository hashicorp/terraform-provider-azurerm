package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudAppId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func SpringCloudAppID(input string) (*SpringCloudAppId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Spring Cloud App ID %q: %+v", input, err)
	}

	app := SpringCloudAppId{
		ResourceGroup: id.ResourceGroup,
	}

	if app.ServiceName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}

	if app.Name, err = id.PopSegment("apps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &app, nil
}
