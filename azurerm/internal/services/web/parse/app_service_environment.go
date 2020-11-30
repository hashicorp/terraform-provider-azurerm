package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceEnvironmentId struct {
	ResourceGroup          string
	HostingEnvironmentName string
}

func AppServiceEnvironmentID(input string) (*AppServiceEnvironmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Environment ID %q: %+v", input, err)
	}

	appServiceEnvironment := AppServiceEnvironmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if appServiceEnvironment.HostingEnvironmentName, err = id.PopSegment("hostingEnvironments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceEnvironment, nil
}
