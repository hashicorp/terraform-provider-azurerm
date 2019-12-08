package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServicePlanResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseAppServicePlanID(input string) (*AppServicePlanResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", input, err)
	}

	group := AppServicePlanResourceID{
		Base: *id,
		Name: id.Path["serverfarms"],
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `serverfarms` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "serverfarms")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
