package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubConnectionId struct {
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func NewVirtualHubConnectionID(id VirtualHubId, name string) VirtualHubConnectionId {
	return VirtualHubConnectionId{
		ResourceGroup:  id.ResourceGroup,
		VirtualHubName: id.Name,
		Name:           name,
	}
}

func (id VirtualHubConnectionId) ID(subscriptionId string) string {
	base := NewVirtualHubID(id.ResourceGroup, id.VirtualHubName).ID(subscriptionId)
	return fmt.Sprintf("%s/hubVirtualNetworkConnections/%s", base, id.Name)
}

func VirtualHubConnectionID(input string) (*VirtualHubConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Hub Connection ID %q: %+v", input, err)
	}

	// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tom-dev99/providers/Microsoft.Network/virtualHubs/tom-devvh/hubVirtualNetworkConnections/first
	connection := VirtualHubConnectionId{
		ResourceGroup:  id.ResourceGroup,
		VirtualHubName: id.Path["virtualHubs"],
		Name:           id.Path["hubVirtualNetworkConnections"],
	}

	if connection.VirtualHubName == "" {
		return nil, fmt.Errorf("ID was missing the `virtualHubs` element")
	}

	if connection.Name == "" {
		return nil, fmt.Errorf("ID was missing the `hubVirtualNetworkConnections` element")
	}

	return &connection, nil
}
