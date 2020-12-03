package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridConnectionId struct {
	ResourceGroup string
	Name          string
	NamespaceName string
}

func HybridConnectionID(input string) (*HybridConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Hybrid Connection ID %q: %+v", input, err)
	}
	hybridConnection := HybridConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if hybridConnection.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}

	if hybridConnection.Name, err = id.PopSegment("hybridConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &hybridConnection, nil
}
