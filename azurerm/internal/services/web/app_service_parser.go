package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseAppServiceID(input string) (*AppServiceResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", input, err)
	}

	group := AppServiceResourceID{
		Base: *id,
		Name: id.Path["sites"],
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `sites` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "sites")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
