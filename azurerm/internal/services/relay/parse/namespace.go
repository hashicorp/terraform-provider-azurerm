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
		return nil, fmt.Errorf("[ERROR] Unable to parse Relay Namespace ID %q: %+v", input, err)
	}
	nameSpace := NamespaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if nameSpace.Name, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &nameSpace, nil
}
