package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridConnectionId struct {
	ResourceGroup                 string
	SiteName                      string
	RelayName                     string
	HybridConnectionNamespaceName string
}

func HybridConnectionID(input string) (*HybridConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Hybrid Connection ID %q: %+v", input, err)
	}

	hybridConnection := HybridConnectionId{
		ResourceGroup: id.ResourceGroup,
	}
	if hybridConnection.RelayName, err = id.PopSegment("relays"); err != nil {
		return nil, err
	}

	if hybridConnection.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if hybridConnection.HybridConnectionNamespaceName, err = id.PopSegment("hybridConnectionNamespaces"); err != nil {
		return nil, err
	}

	return &hybridConnection, nil
}
