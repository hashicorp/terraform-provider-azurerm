package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceBusNamespaceId struct {
	Name          string
	ResourceGroup string
}

func ServiceBusNamespaceID(input string) (*ServiceBusNamespaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace ID %q: %+v", input, err)
	}

	namespace := ServiceBusNamespaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if namespace.Name, err = id.PopSegment("namespaces"); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace ID %q: %+v", input, err)
	}

	return &namespace, nil
}
