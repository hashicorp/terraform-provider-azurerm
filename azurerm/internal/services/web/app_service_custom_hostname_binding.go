package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCustomHostnameBindingResourceID struct {
	ResourceGroup  string
	AppServiceName string
	Name           string
}

func ParseAppServiceCustomHostnameBindingID(input string) (*AppServiceCustomHostnameBindingResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Custom Hostname Binding ID %q: %+v", input, err)
	}

	group := AppServiceCustomHostnameBindingResourceID{
		ResourceGroup: id.ResourceGroup,
	}
	group.AppServiceName, err = id.PopSegment("sites")
	if err != nil {
		return nil, err
	}

	group.Name, err = id.PopSegment("hostNameBindings")
	if err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
