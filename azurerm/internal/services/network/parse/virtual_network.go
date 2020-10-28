package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkId struct {
	ResourceGroup string
	Name          string
}

func (id VirtualNetworkId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualNetworkID(resourceGroup, name string) VirtualNetworkId {
	return VirtualNetworkId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func VirtualNetworkID(input string) (*VirtualNetworkId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Network ID %q: %+v", input, err)
	}

	vnet := VirtualNetworkId{
		ResourceGroup: id.ResourceGroup,
	}

	if vnet.Name, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vnet, nil
}
