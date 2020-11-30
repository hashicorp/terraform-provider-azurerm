package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CertificateOrderId struct {
	ResourceGroup string
	Name          string
}

func CertificateOrderID(input string) (*CertificateOrderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate Order ID %q: %+v", input, err)
	}

	order := CertificateOrderId{
		ResourceGroup: id.ResourceGroup,
	}

	if order.Name, err = id.PopSegment("certificateOrders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &order, nil
}
