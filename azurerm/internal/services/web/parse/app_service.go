package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceResourceID struct {
	ResourceGroup string
	Name          string
}

func AppServiceID(input string) (*AppServiceResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service ID %q: %+v", input, err)
	}

	appService := AppServiceResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if appService.Name, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appService, nil
}
