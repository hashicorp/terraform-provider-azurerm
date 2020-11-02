package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceFabricMeshNetworkId struct {
	ResourceGroup string
	Name          string
}

func ServiceFabricMeshNetworkID(input string) (*ServiceFabricMeshNetworkId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Service Fabric Mesh Network ID %q: %+v", input, err)
	}

	network := ServiceFabricMeshNetworkId{
		ResourceGroup: id.ResourceGroup,
	}

	if network.Name, err = id.PopSegment("networks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &network, nil
}
