package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCustomHostnameBindingId struct {
	ResourceGroup  string
	AppServiceName string
	Name           string
}

func AppServiceCustomHostnameBindingID(input string) (*AppServiceCustomHostnameBindingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Custom Hostname Binding ID %q: %+v", input, err)
	}

	binding := AppServiceCustomHostnameBindingId{
		ResourceGroup: id.ResourceGroup,
	}

	if binding.AppServiceName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if binding.Name, err = id.PopSegment("hostNameBindings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &binding, nil
}
