package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NamespaceId struct {
	ResourceGroup string
	Name          string
}

func NamespaceID(input string) (*NamespaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Service Bus Namespace ID %q: %+v", input, err)
	}

	namespace := NamespaceId{
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
