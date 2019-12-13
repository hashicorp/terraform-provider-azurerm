package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCustomHostnameBindingResourceID struct {
	Base azure.ResourceID

	AppServiceName string
	Name           string
}

func ParseAppServiceCustomHostnameBindingID(input string) (*AppServiceCustomHostnameBindingResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Custom Hostname Binding ID %q: %+v", input, err)
	}

	group := AppServiceCustomHostnameBindingResourceID{
		Base:           *id,
		AppServiceName: id.Path["sites"],
		Name:           id.Path["hostNameBindings"],
	}

	if group.AppServiceName == "" {
		return nil, fmt.Errorf("ID was missing the `sites` element")
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `hostNameBindings` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "sites")
	delete(pathWithoutElements, "hostNameBindings")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
