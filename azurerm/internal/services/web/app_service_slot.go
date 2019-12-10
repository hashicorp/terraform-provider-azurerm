package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceSlotResourceID struct {
	Base azure.ResourceID

	AppServiceName string
	Name           string
}

func ParseAppServiceSlotID(input string) (*AppServiceSlotResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Slot ID %q: %+v", input, err)
	}

	group := AppServiceSlotResourceID{
		Base:           *id,
		AppServiceName: id.Path["sites"],
		Name:           id.Path["slots"],
	}

	if group.AppServiceName == "" {
		return nil, fmt.Errorf("ID was missing the `sites` element")
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `slots` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "sites")
	delete(pathWithoutElements, "slots")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
