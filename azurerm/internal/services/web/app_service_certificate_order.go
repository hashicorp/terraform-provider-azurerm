package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCertificateOrderResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseAppServiceCertificateOrderID(input string) (*AppServiceCertificateOrderResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate Order ID %q: %+v", input, err)
	}

	order := AppServiceCertificateOrderResourceID{
		ResourceGroup: id.ResourceGroup,
		Name:          id.Path[""],
	}
	order.Name, err = id.PopSegment("certificateOrders")
	if err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &order, nil
}
