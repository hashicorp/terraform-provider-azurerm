package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceEnvironmentResourceID struct {
	ResourceGroup string
	Name          string
}

func AppServiceEnvironmentID(input string) (*AppServiceEnvironmentResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Environment ID %q: %+v", input, err)
	}

	appServiceEnvironment := AppServiceEnvironmentResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if appServiceEnvironment.Name, err = id.PopSegment("hostingEnvironments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceEnvironment, nil
}
