package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCertificateOrderResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseAppServiceCertificateOrderID(input string) (*AppServiceCertificateOrderResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate Order ID %q: %+v", input, err)
	}

	group := AppServiceCertificateOrderResourceID{
		Base: *id,
		Name: id.Path["certificateOrders"],
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `certificateOrders` element")
	}

	pathWithoutElements := group.Base.Path
	delete(pathWithoutElements, "certificateOrders")
	if len(pathWithoutElements) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
